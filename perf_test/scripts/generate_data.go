package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

var (
	offerTypes    = []string{"sale", "rent"}
	propertyTypes = []string{"house", "apartment"}
	statuses      = []string{"active", "sold", "archived"}
	statusWeights = []float32{0.7, 0.2, 0.1} // 70% active, 20% sold, 10% archived
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run generate_data.go <dsn>")
	}

	dsn := os.Args[1]
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rand.Seed(time.Now().UnixNano())

	ctx := context.Background()

	log.Println("🚀 Starting data generation...")

	// Создаем тестовых пользователей
	log.Println("📝 Creating users...")
	userIDs := createUsers(ctx, db, 1000)
	log.Printf("✅ Created %d users", len(userIDs))

	// Создаем регионы и локации
	log.Println("📍 Creating locations...")
	locationIDs := createLocations(ctx, db, 200)
	log.Printf("✅ Created %d locations", len(locationIDs))

	// Создаем станции метро
	log.Println("🚇 Creating metro stations...")
	metroIDs := createMetroStations(ctx, db, locationIDs, 50)
	log.Printf("✅ Created %d metro stations", len(metroIDs))

	// Связываем локации с метро
	log.Println("🔗 Linking locations with metro...")
	createLocationMetroLinks(ctx, db, locationIDs, metroIDs)

	// Создаем жилые комплексы
	log.Println("🏢 Creating housing complexes...")
	complexIDs := createComplexes(ctx, db, locationIDs, 100)
	log.Printf("✅ Created %d housing complexes", len(complexIDs))

	// Создаем фото для комплексов
	log.Println("📷 Creating complex photos...")
	createComplexPhotos(ctx, db, complexIDs)

	// Создаем объявления
	log.Println("🏠 Creating offers (this may take a while)...")
	offerIDs := createOffers(ctx, db, userIDs, locationIDs, complexIDs, 50000)
	log.Printf("✅ Created %d offers", len(offerIDs))

	// Создаем фото для объявлений
	log.Println("📸 Creating offer photos...")
	createOfferPhotos(ctx, db, offerIDs)

	// Создаем профили для пользователей
	log.Println("👤 Creating user profiles...")
	createProfiles(ctx, db, userIDs)

	log.Println("🎉 Data generation completed successfully!")
}

func createUsers(ctx context.Context, db *sql.DB, count int) []uuid.UUID {
	ids := make([]uuid.UUID, 0, count)
	roles := []string{"user", "owner", "realtor"}
	roleWeights := []float32{0.6, 0.3, 0.1}

	// Создаем тестовый хеш пароля один раз
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("TestPassword123!"), bcrypt.MinCost)

	stmt, err := db.PrepareContext(ctx,
		`INSERT INTO users (id, email, password_hash, role) 
		 VALUES ($1, $2, $3, $4)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; i < count; i++ {
		id := uuid.New()
		email := fmt.Sprintf("loadtest_user_%d@example.com", i)
		role := weightedChoice(roles, roleWeights)

		_, err := stmt.ExecContext(ctx, id, email, string(passwordHash), role)
		if err != nil {
			log.Printf("Error creating user %d: %v", i, err)
			continue
		}
		ids = append(ids, id)

		if i%100 == 0 && i > 0 {
			log.Printf("  Created %d/%d users", i, count)
		}
	}

	return ids
}

func createProfiles(ctx context.Context, db *sql.DB, userIDs []uuid.UUID) {
	firstNames := []string{"Иван", "Мария", "Алексей", "Елена", "Дмитрий", "Анна", "Сергей", "Ольга"}
	lastNames := []string{"Иванов", "Петров", "Сидоров", "Смирнов", "Кузнецов", "Попов", "Васильев"}

	stmt, err := db.PrepareContext(ctx,
		`INSERT INTO profile (id, user_id, first_name, last_name, phone) 
		 VALUES ($1, $2, $3, $4, $5)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i, userID := range userIDs {
		id := uuid.New()
		firstName := firstNames[rand.Intn(len(firstNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]
		phone := fmt.Sprintf("+7-9%02d-%03d-%02d-%02d",
			rand.Intn(100), rand.Intn(1000), rand.Intn(100), rand.Intn(100))

		_, err := stmt.ExecContext(ctx, id, userID, firstName, lastName, phone)
		if err != nil {
			log.Printf("Error creating profile for user %d: %v", i, err)
		}
	}
}

func createLocations(ctx context.Context, db *sql.DB, count int) []uuid.UUID {
	ids := make([]uuid.UUID, 0, count)

	// Создаем регионы
	regions := []struct {
		name string
		slug string
	}{
		{"Москва", "moscow"},
		{"Санкт-Петербург", "saint-petersburg"},
		{"Московская область", "moscow-region"},
	}

	regionIDs := make([]uuid.UUID, 0, len(regions))
	for _, r := range regions {
		id := uuid.New()
		_, err := db.ExecContext(ctx,
			`INSERT INTO region (id, name, level, slug) 
			 VALUES ($1, $2, 0, $3)`,
			id, r.name, r.slug)
		if err != nil {
			log.Printf("Error creating region %s: %v", r.name, err)
			continue
		}
		regionIDs = append(regionIDs, id)
	}

	// Создаем локации
	stmt, err := db.PrepareContext(ctx,
		`INSERT INTO location (id, region_id, latitude, longitude) 
		 VALUES ($1, $2, $3, $4)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; i < count; i++ {
		id := uuid.New()
		regionID := regionIDs[rand.Intn(len(regionIDs))]

		// Координаты Москвы с разбросом
		lat := 55.7558 + (rand.Float64()-0.5)*0.5
		lon := 37.6173 + (rand.Float64()-0.5)*0.5

		_, err := stmt.ExecContext(ctx, id, regionID, lat, lon)
		if err != nil {
			log.Printf("Error creating location %d: %v", i, err)
			continue
		}
		ids = append(ids, id)
	}

	return ids
}

func createMetroStations(ctx context.Context, db *sql.DB, locationIDs []uuid.UUID, count int) []uuid.UUID {
	ids := make([]uuid.UUID, 0, count)
	stationNames := []string{
		"Красные Ворота", "Парк Культуры", "Кропоткинская", "Арбатская",
		"Смоленская", "Киевская", "Маяковская", "Белорусская", "Новослободская",
		"Проспект Мира", "Комсомольская", "Курская", "Таганская", "Павелецкая",
		"Добрынинская", "Октябрьская", "Тульская", "Нагатинская", "Коломенская",
	}

	stmt, err := db.PrepareContext(ctx,
		`INSERT INTO metro_station (id, name, location_id) 
		 VALUES ($1, $2, $3)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; i < count && i < len(stationNames)*3; i++ {
		id := uuid.New()
		name := stationNames[i%len(stationNames)]
		if i >= len(stationNames) {
			name = fmt.Sprintf("%s-%d", name, i/len(stationNames))
		}
		locationID := locationIDs[rand.Intn(len(locationIDs))]

		_, err := stmt.ExecContext(ctx, id, name, locationID)
		if err != nil {
			log.Printf("Error creating metro station %d: %v", i, err)
			continue
		}
		ids = append(ids, id)
	}

	return ids
}

func createLocationMetroLinks(ctx context.Context, db *sql.DB, locationIDs, metroIDs []uuid.UUID) {
	stmt, err := db.PrepareContext(ctx,
		`INSERT INTO location_metro (location_id, metro_station_id, distance_meters) 
		 VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Каждая локация связана с 1-3 станциями метро
	for _, locationID := range locationIDs {
		numStations := rand.Intn(3) + 1
		usedStations := make(map[uuid.UUID]bool)

		for j := 0; j < numStations; j++ {
			metroID := metroIDs[rand.Intn(len(metroIDs))]
			if usedStations[metroID] {
				continue
			}
			usedStations[metroID] = true

			distance := rand.Intn(1500) + 100 // 100-1600 метров

			_, err := stmt.ExecContext(ctx, locationID, metroID, distance)
			if err != nil {
				log.Printf("Error linking location to metro: %v", err)
			}
		}
	}
}

func createComplexes(ctx context.Context, db *sql.DB, locationIDs []uuid.UUID, count int) []uuid.UUID {
	ids := make([]uuid.UUID, 0, count)
	developers := []string{
		"ПИК", "Самолет", "ЛСР", "Главстрой", "А101", "Capital Group",
		"МИЦ", "ФСК", "Донстрой", "Vesper",
	}

	stmt, err := db.PrepareContext(ctx,
		`INSERT INTO housing_complex 
		 (id, name, description, year_built, location_id, developer, address, starting_price) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := range count {
		id := uuid.New()
		locationID := locationIDs[rand.Intn(len(locationIDs))]
		name := fmt.Sprintf("ЖК %s-%d", developers[rand.Intn(len(developers))], i)
		description := fmt.Sprintf("Современный жилой комплекс с развитой инфраструктурой. "+
			"Квартиры от застройщика. Комплекс номер %d", i)
		yearBuilt := 2015 + rand.Intn(10)
		developer := developers[rand.Intn(len(developers))]
		address := fmt.Sprintf("ул. Тестовая, д. %d", rand.Intn(100)+1)
		startingPrice := int64(3000000 + rand.Intn(50000000))

		_, err := stmt.ExecContext(ctx, id, name, description, yearBuilt,
			locationID, developer, address, startingPrice)
		if err != nil {
			log.Printf("Error creating complex %d: %v", i, err)
			continue
		}
		ids = append(ids, id)

		if i%20 == 0 && i > 0 {
			log.Printf("  Created %d/%d complexes", i, count)
		}
	}

	return ids
}

func createComplexPhotos(ctx context.Context, db *sql.DB, complexIDs []uuid.UUID) {
	stmt, err := db.PrepareContext(ctx,
		`INSERT INTO complex_photo (id, complex_id, url) 
		 VALUES ($1, $2, $3)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, complexID := range complexIDs {
		// 2-5 фото на комплекс
		numPhotos := rand.Intn(4) + 2
		for j := 0; j < numPhotos; j++ {
			id := uuid.New()
			url := fmt.Sprintf("https://example.com/complex/%s/photo_%d.jpg", complexID, j)

			_, err := stmt.ExecContext(ctx, id, complexID, url)
			if err != nil {
				log.Printf("Error creating complex photo: %v", err)
			}
		}
	}
}

func createOffers(ctx context.Context, db *sql.DB, userIDs, locationIDs, complexIDs []uuid.UUID, count int) []uuid.UUID {
	ids := make([]uuid.UUID, 0, count)

	stmt, err := db.PrepareContext(ctx,
		`INSERT INTO offer 
		 (id, user_id, location_id, housing_complex_id, title, description, 
		  price, area, address, rooms, property_type, offer_type, status, 
		  floor, total_floors, living_area, kitchen_area, deposit, commission, rental_period) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; i < count; i++ {
		id := uuid.New()
		userID := userIDs[rand.Intn(len(userIDs))]
		locationID := locationIDs[rand.Intn(len(locationIDs))]

		var complexID *uuid.UUID
		if rand.Float32() < 0.6 { // 60% объявлений в комплексах
			cid := complexIDs[rand.Intn(len(complexIDs))]
			complexID = &cid
		}

		offerType := offerTypes[rand.Intn(len(offerTypes))]
		propertyType := propertyTypes[rand.Intn(len(propertyTypes))]
		status := weightedChoice(statuses, statusWeights)
		rooms := rand.Intn(5) + 1

		title := generateTitle(propertyType, offerType, rooms)
		description := generateDescription(propertyType, rooms)

		var price int64
		if offerType == "sale" {
			price = int64(2000000 + rand.Intn(50000000))
		} else {
			price = int64(20000 + rand.Intn(150000))
		}

		area := float64(30 + rand.Intn(170))
		livingArea := area * (0.5 + rand.Float64()*0.3)
		kitchenArea := 8.0 + rand.Float64()*12.0

		floor := rand.Intn(25) + 1
		totalFloors := floor + rand.Intn(10)
		address := fmt.Sprintf("ул. Тестовая, д. %d, кв. %d", rand.Intn(100)+1, rand.Intn(200)+1)

		var deposit, commission *int64
		var rentalPeriod *string
		if offerType == "rent" {
			d := price // депозит = месячная аренда
			c := int64(float64(price) * 0.5)
			rp := "длительный срок"
			deposit = &d
			commission = &c
			rentalPeriod = &rp
		}

		_, err := stmt.ExecContext(ctx, id, userID, locationID, complexID,
			title, description, price, area, address, rooms, propertyType,
			offerType, status, floor, totalFloors, livingArea, kitchenArea,
			deposit, commission, rentalPeriod)

		if err != nil {
			log.Printf("Error creating offer %d: %v", i, err)
			continue
		}
		ids = append(ids, id)

		if i%1000 == 0 && i > 0 {
			log.Printf("  Created %d/%d offers", i, count)
		}
	}

	return ids
}

func createOfferPhotos(ctx context.Context, db *sql.DB, offerIDs []uuid.UUID) {
	stmt, err := db.PrepareContext(ctx,
		`INSERT INTO offer_photo (id, offer_id, url) 
		 VALUES ($1, $2, $3)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i, offerID := range offerIDs {
		// 3-10 фото на объявление
		numPhotos := rand.Intn(8) + 3
		for j := 0; j < numPhotos; j++ {
			id := uuid.New()
			url := fmt.Sprintf("https://example.com/offer/%s/photo_%d.jpg", offerID, j)

			_, err := stmt.ExecContext(ctx, id, offerID, url)
			if err != nil {
				log.Printf("Error creating offer photo: %v", err)
			}
		}

		if i%1000 == 0 && i > 0 {
			log.Printf("  Created photos for %d/%d offers", i, len(offerIDs))
		}
	}
}

// Helper functions

func weightedChoice(items []string, weights []float32) string {
	totalWeight := float32(0)
	for _, w := range weights {
		totalWeight += w
	}

	r := rand.Float32() * totalWeight
	cumWeight := float32(0)

	for i, w := range weights {
		cumWeight += w
		if r < cumWeight {
			return items[i]
		}
	}

	return items[len(items)-1]
}

func generateTitle(propertyType, offerType string, rooms int) string {
	action := "Продажа"
	if offerType == "rent" {
		action = "Аренда"
	}

	propType := "квартиры"
	if propertyType == "house" {
		propType = "дома"
	}

	if rooms == 0 {
		return fmt.Sprintf("%s студии", action)
	}

	return fmt.Sprintf("%s %d-комнатной %s", action, rooms, propType)
}

func generateDescription(propertyType string, rooms int) string {
	descriptions := []string{
		"Отличное состояние, свежий ремонт. Все коммуникации в порядке.",
		"Требуется косметический ремонт. Хорошая планировка.",
		"Евроремонт, встроенная техника. Готово к заселению.",
		"Студия в современном доме. Панорамные окна.",
		"Просторная квартира с большой кухней. Тихий район.",
	}

	return descriptions[rand.Intn(len(descriptions))]
}
