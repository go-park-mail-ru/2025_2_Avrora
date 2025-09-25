# Схема базы данных: Объявления недвижимости

Схема состоит из шести отношений: `User`, `Offer`, `Location`, `Photo`, `Category`, `Region`.  
Все отношения находятся в **Нормальной форме Бойса-Кодда (НФБК)**.  

---

##  Отношение: User

**Описание**:  
Хранит уникальных пользователей системы. Каждый пользователь имеет email и хеш пароля.

### Функциональные зависимости:

{ID} → Email, Password  
{Email} → ID, Password

> **Пояснение**:  
> - `ID` — первичный ключ, однозначно определяет запись
> - `Email` — уникальный атрибут и тоже суперключ

---

##  Отношение: Category

**Описание**:  
Классифицирует типы недвижимости: квартира, дом, комната и тд + легче фильтровать

### Функциональные зависимости:

{ID} → Name, Slug, Description  
{Slug} → ID, Name, Description

> **Пояснение**:  
> - `ID` — первичный ключ.  
> - `Slug` — уникальный идентификатор для юзера. пример: 'kvartira'

---

##  Отношение: Region

**Описание**:  
Иерархическая таблица регионов: Страна → Регион → Город → Район. Упрощает поиск и аналитику.

### Функциональные зависимости:

{ID} → Name, ParentID, Level, Slug  
{Slug} → ID, Name, ParentID, Level  
{ParentID, Name} → ID, Slug, Level

> **Пояснение**:  
> - `ID` — первичный ключ.  
> - `Slug` — уникальный → суперключ.  
> - `{ParentID, Name}` — уникальная комбинация внутри родителя → суперключ.  
> - `ParentID` ссылается на `Region.ID` → рекурсивная связь.

---

## Отношение: Location

**Описание**:  
Хранит нормализованные локации (адреса) для объявлений. Позволяет избежать дублирования одинаковых адресов.

### Функциональные зависимости:

{ID} → RegionID, Street, HouseNumber, Latitude, Longitude  
{RegionID, Street, HouseNumber} → ID, Latitude, Longitude

> **Пояснение**:  
> - `ID` — суррогатный первичный ключ.  
> - Комбинация `{RegionID, Street, HouseNumber}` — уникальный ключ

---

## Отношение: Offer

**Описание**:  
Хранит объявления о недвижимости. Каждое объявление принадлежит ровно одному пользователю, категории и локации.

### Функциональные зависимости:

{ID} → UserID, LocationID, CategoryID, Title, Description, Price, Area, Rooms, OfferType, CreatedAt, UpdatedAt  
{UserID, Title, LocationID, CreatedAt} → ID, Description, Price, Area, Rooms, OfferType, UpdatedAt

> **Пояснение**:  
> - `ID` — первичный ключ.  
> - Комбинация `{UserID, Title, LocationID, CreatedAt}` — уникальный бизнес-ключ → суперключ.  
> - Адрес и категория вынесены в отдельные таблицы → устранена избыточность.

---

##  Отношение: Photo

**Описание**:  
Хранит ссылки на фотографии для объявлений. Одно объявление может иметь множество фотографий.

### Функциональные зависимости:

{ID} → OfferID, URL, Position, UploadedAt  
{OfferID, Position} → ID, URL, UploadedAt  
{OfferID, URL} → ID, Position, UploadedAt

> **Пояснение**:  
> - `ID` — первичный ключ.  
> - `{OfferID, Position}` — уникальная позиция фото в галерее → суперключ.  
> - `{OfferID, URL}` — одна и та же ссылка не может быть прикреплена дважды к одному объявлению → суперключ.

---

## ERD

```mermaid
erDiagram

    direction TB

    User ||--o{ Offer : "создаёт"
    Category ||--o{ Offer : "категория"
    Region ||--o{ Location : "содержит"
    Location ||--o{ Offer : "расположен в"
    Offer ||--o{ Photo : "имеет"

    User {
        int ID PK
        string Email UK
        string Password
    }

    Category {
        int ID PK
        string Name
        string Slug UK
        string Description
    }

    Region {
        int ID PK
        string Name
        int ParentID FK "может быть NULL"
        int Level
        string Slug UK
    }

    Location {
        int ID PK
        int RegionID FK
        string Street
        string HouseNumber
        decimal Latitude
        decimal Longitude
    }

    Offer {
        int ID PK
        int UserID FK
        int LocationID FK
        int CategoryID FK
        string Title
        string Description
        bigint Price
        decimal Area
        int Rooms
        string OfferType
        timestamp CreatedAt
        timestamp UpdatedAt
    }

    Photo {
        int ID PK
        int OfferID FK
        string URL
        int Position
        timestamp UploadedAt
    }