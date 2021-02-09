package main

import (
	"fmt"
	"os"

	"github.com/HaiderAliHosen/sqlitedemo/controllers"
	"github.com/HaiderAliHosen/sqlitedemo/infrastructure"
	"github.com/HaiderAliHosen/sqlitedemo/middlewares"
	"github.com/HaiderAliHosen/sqlitedemo/models"
	"github.com/HaiderAliHosen/sqlitedemo/seeds"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func drop(db *gorm.DB) {
	db.DropTableIfExists(
		&models.FileUpload{},
		&models.Comment{},
		&models.OrderItem{}, &models.Order{}, &models.Address{},
		&models.ProductCategory{}, &models.ProductTag{},
		&models.Tag{}, &models.Category{},
		&models.Product{},
		&models.UserRole{}, &models.Role{}, &models.User{})
}

func migrate(database *gorm.DB) {

	database.AutoMigrate(&models.Address{})

	database.AutoMigrate(&models.Category{})
	database.AutoMigrate(&models.Comment{})

	database.AutoMigrate(&models.Order{})
	database.AutoMigrate(&models.OrderItem{})

	database.AutoMigrate(&models.Product{})
	database.AutoMigrate(&models.ProductCategory{})

	database.AutoMigrate(&models.Role{})
	database.AutoMigrate(&models.UserRole{})

	database.AutoMigrate(&models.Tag{})
	database.AutoMigrate(&models.ProductTag{})

	database.AutoMigrate(&models.User{})

	database.AutoMigrate(&models.FileUpload{})
}

func addDbConstraints(database *gorm.DB) {
	// TODO: it is well known GORM does not add foreign keys even after using ForeignKey in struct, but, why manually does not work neither ?

	dialect := database.Dialect().GetName() // mysql, sqlite3
	if dialect != "sqlite3" {
		database.Model(&models.Comment{}).AddForeignKey("product_id", "products(id)", "CASCADE", "CASCADE")
		database.Model(&models.Comment{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

		database.Model(&models.Order{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
		database.Model(&models.Order{}).AddForeignKey("address_id", "addresses(id)", "CASCADE", "CASCADE")
		database.Model(&models.OrderItem{}).AddForeignKey("order_id", "orders(id)", "CASCADE", "CASCADE")
		database.Model(&models.OrderItem{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

		database.Model(&models.Address{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

		database.Model(&models.UserRole{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
		database.Model(&models.UserRole{}).AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE")

		database.Table("products_tags").AddForeignKey("product_id", "products(id)", "CASCADE", "CASCADE")
		database.Table("products_tags").AddForeignKey("tag_id", "tags(id)", "CASCADE", "CASCADE")

		database.Model(models.ProductCategory{}).AddForeignKey("product_id", "products(id)", "CASCADE", "CASCADE")
		database.Model(models.ProductCategory{}).AddForeignKey("category_id", "categories(id)", "CASCADE", "CASCADE")
	} else if dialect == "sqlite3" {
		database.Table("comments").AddIndex("comments__idx_product_id", "product_id")
		database.Table("comments").AddIndex("comments__idx_user_id", "user_id")

		database.Table("ratings").AddIndex("ratings__idx_user_id", "user_id")
		database.Table("ratings").AddIndex("ratings__idx_product_id", "product_id")

		database.Model(&models.Comment{}).AddIndex("comments__idx_created_at", "created_at")

	}

	database.Model(&models.UserRole{}).AddIndex("user_roles__idx_user_id", "user_id")
	database.Table("products_tags").AddIndex("products_tags__idx_product_id", "product_id")
}
func create(database *gorm.DB) {
	drop(database)
	migrate(database)
	addDbConstraints(database)
}

func main() {

	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}
	println(os.Getenv("DB_DIALECT"))

	database := infrastructure.OpenDbConnection()

	defer database.Close()
	args := os.Args
	if len(args) > 1 {
		first := args[1]
		second := ""
		if len(args) > 2 {
			second = args[2]
		}

		if first == "create" {
			create(database)
		} else if first == "seed" {
			seeds.Seed()
			os.Exit(0)
		} else if first == "migrate" {
			migrate(database)
		}

		if second == "seed" {
			seeds.Seed()
			os.Exit(0)
		} else if first == "migrate" {
			migrate(database)
		}

		if first != "" && second == "" {
			os.Exit(0)
		}
	}

	migrate(database)

	// gin.New() - new gin Instance with no middlewares
	// goGonicEngine.Use(gin.Logger())
	// goGonicEngine.Use(gin.Recovery())
	goGonicEngine := gin.Default() // gin with the Logger and Recovery Middlewares attached
	//goGonicEngine.Use(cors.Default())
	// Allow all Origins
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("authorization")
	goGonicEngine.Use(cors.New(config))

	// goGonicEngine.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"*"},
	// 	AllowMethods:     []string{"POST, OPTIONS, GET, PUT"},
	// 	AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	// AllowOriginFunc: func(origin string) bool {
	// 	// 	return origin == "https://github.com"
	// 	// },
	// 	MaxAge: 12 * time.Hour,
	// }))

	// goGonicEngine.Use(cors.New(cors.Config{
	// 	AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT"},
	// 	AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	AllowAllOrigins:  true,
	// 	AllowOriginFunc:  func(origin string) bool { return true },
	// 	MaxAge:           86400,
	// }))

	//goGonicEngine.Use(cors.Default())

	goGonicEngine.Use(middlewares.Benchmark())

	//goGonicEngine.Use(middlewares.Cors())
	//goGonicEngine.Use(middlewares.CORSMiddleware())

	//goGonicEngine.Use(middlewares.UserLoaderMiddleware())
	goGonicEngine.Static("/static", "./static")
	apiRouteGroup := goGonicEngine.Group("/api")

	controllers.RegisterUserRoutes(apiRouteGroup.Group("/users"))
	controllers.RegisterProductRoutes(apiRouteGroup.Group("/products", middlewares.UserLoaderMiddleware()))
	controllers.RegisterCommentRoutes(apiRouteGroup.Group("/", middlewares.UserLoaderMiddleware()))
	controllers.RegisterPageRoutes(apiRouteGroup.Group("/", middlewares.UserLoaderMiddleware()))
	controllers.RegisterAddressesRoutes(apiRouteGroup.Group("/users", middlewares.UserLoaderMiddleware()))
	controllers.RegisterTagRoutes(apiRouteGroup.Group("/tags", middlewares.UserLoaderMiddleware()))
	controllers.RegisterCategoryRoutes(apiRouteGroup.Group("/categories", middlewares.UserLoaderMiddleware()))
	controllers.RegisterOrderRoutes(apiRouteGroup.Group("/orders", middlewares.UserLoaderMiddleware()))

	goGonicEngine.Run(":8080") // listen and serve on 0.0.0.0:8080
}
