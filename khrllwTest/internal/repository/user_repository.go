package repository

import (
	"gorm.io/gorm"
	"khrllwTest/internal/models"
)

// UserRepositoryImpl - реализация для GORM
type UserRepositoryImpl struct {
	db *gorm.DB // Экземпляр подключения к БД
}

// UserRepository определяет контракт для работы с пользователями в БД
type UserRepository interface {
	// Create
	//Создание нового пользователя
	Create(user *models.User) error

	// FindByID
	//Поиск пользователя по ID (возвращает nil, если не найден)
	FindByID(id uint) (*models.User, error)

	// FindByEmail
	//Поиск пользователя по email (возвращает nil, если не найден)
	FindByEmail(email string) (*models.User, error)

	// Update
	//Обновление данных пользователя
	Update(user *models.User) error

	// Delete
	//Удаление пользователя по ID (мягкое удаление, если настроено в модели)
	Delete(id uint) error

	// GetAll
	// Получение списка пользователей с пагинацией и фильтрацией
	GetAll(offset, limit, minAge, maxAge int) ([]models.User, int64, error)
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(user *models.User) error {
	// Выполняет INSERT запрос
	return r.db.Create(user).Error
}

func (r *UserRepositoryImpl) FindByID(id uint) (*models.User, error) {
	var user models.User
	// SELECT * FROM users WHERE id = ? LIMIT 1
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*models.User, error) {
	var user models.User
	// SELECT * FROM users WHERE email = ? LIMIT 1
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepositoryImpl) Update(user *models.User) error {
	// Выполняет UPDATE запрос
	return r.db.Save(user).Error
}

func (r *UserRepositoryImpl) Delete(id uint) error {
	// DELETE FROM users WHERE id = ?
	return r.db.Delete(&models.User{}, id).Error
}

// GetAll возвращает список пользователей с пагинацией и фильтрацией по возрасту
func (r *UserRepositoryImpl) GetAll(offset, limit, minAge, maxAge int) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	query := r.db.Model(&models.User{})
	// Добавление фильтрации
	if minAge > 0 {
		query = query.Where("age >= ?", minAge)
	}
	if maxAge > 0 {
		query = query.Where("age <= ?", maxAge)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// Применение пагинации
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
