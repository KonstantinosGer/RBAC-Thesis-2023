package routes

import (
	"backend/config"
	"backend/handlers"
	"backend/middleware"
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"net/http"
)

//SetupRoutes : all the routes are defined here
func SetupRoutes(db *gorm.DB) {
	httpRouter := gin.Default()

	//CORS
	cors_conf := cors.DefaultConfig()
	cors_conf.AllowAllOrigins = true
	cors_conf.AllowCredentials = true
	cors_conf.AddAllowHeaders("authorization")
	cors_conf.AddAllowHeaders("Access-Control-Allow-Credentials")
	cors_conf.AddAllowHeaders("Access-Control-Allow-Origin")
	cors_conf.AddAllowHeaders("accept")
	httpRouter.Use(cors.New(cors_conf))
	httpRouter.MaxMultipartMemory = 1024 << 20

	//------
	//Casbin
	//------
	// Initialize  casbin adapter
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
	}

	// Load models configuration file and policy store adapter
	enforcer, err := casbin.NewSyncedEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		panic(fmt.Sprintf("failed to create casbin enforcer: %v", err))
	}

	//--------
	//Firebase
	//--------
	// configure firebase
	firebaseAuth := config.SetupFirebase()

	// set db & firebase auth to gin context with a middleware to all incoming request
	httpRouter.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("firebaseAuth", firebaseAuth)
	})

	apiRoutes := httpRouter.Group("/api", middleware.AuthMiddleware)

	//------------
	//USERS ROUTES
	//------------
	userProtectedRoutes := apiRoutes.Group("/users") //, middleware.AuthMiddleware
	{
		userProtectedRoutes.GET("/unassigned", middleware.Authorize("rbac::data", "read", enforcer), handlers.GetUnassignedUsers())
		userProtectedRoutes.GET("/emails", middleware.Authorize("rbac::data", "read", enforcer), handlers.GetUsersEmails())
		userProtectedRoutes.GET("/", middleware.Authorize("rbac::data", "read", enforcer), handlers.GetAllUsers())
		userProtectedRoutes.GET("/sync", handlers.SyncUsersWithFirebase())
	}

	//------------
	//ROLES ROUTES
	//------------
	roleProtectedRoutes := apiRoutes.Group("/roles")
	{
		roleProtectedRoutes.PUT("/", middleware.Authorize("rbac::data", "write", enforcer), handlers.UpdateRole(enforcer))
		roleProtectedRoutes.GET("/", middleware.Authorize("rbac::data", "read", enforcer), handlers.GetAllRoles())
		roleProtectedRoutes.POST("/", middleware.Authorize("rbac::data", "write", enforcer), handlers.AddRole())
		roleProtectedRoutes.DELETE("/", middleware.Authorize("rbac::data", "write", enforcer), handlers.DeleteRole(enforcer))
	}

	//------------
	//PERMISSIONS ROUTES
	//------------
	permissionProtectedRoutes := apiRoutes.Group("/permissions")
	{
		permissionProtectedRoutes.GET("/", middleware.Authorize("rbac::data", "read", enforcer), handlers.GetPermissionsForRole(enforcer))
		permissionProtectedRoutes.POST("/", middleware.Authorize("rbac::data", "write", enforcer), handlers.AddPermission(enforcer))
		permissionProtectedRoutes.DELETE("/", middleware.Authorize("rbac::data", "write", enforcer), handlers.DeletePermission(enforcer))
	}

	//------------
	//CASBIN ROUTES
	//------------
	casbinProtectedRoutes := apiRoutes.Group("/casbin")
	{
		casbinProtectedRoutes.POST("/permissions", handlers.GetFrontendPermission(enforcer))
	}

	//------------
	//EMPLOYEES ROUTES
	//------------
	employeesProtectedRoutes := apiRoutes.Group("/employees")
	{
		employeesProtectedRoutes.GET("/", middleware.Authorize("rbac::data", "read", enforcer), handlers.GetAllEmployees())
		employeesProtectedRoutes.POST("/", middleware.Authorize("rbac::data", "write", enforcer), handlers.AddEmployee())
		employeesProtectedRoutes.DELETE("/:id", middleware.Authorize("rbac::data", "write", enforcer), handlers.DeleteEmployee(enforcer))
		employeesProtectedRoutes.PUT("/", middleware.Authorize("rbac::data", "write", enforcer), handlers.UpdateEmployee())

		associations := employeesProtectedRoutes.Group("/associations")
		{
			associations.POST("/", middleware.Authorize("rbac::data", "write", enforcer), handlers.AddAssociation())
			associations.DELETE("/", middleware.Authorize("rbac::data", "write", enforcer), handlers.DeleteAssociation())
			associations.GET("/:id", middleware.Authorize("rbac::data", "read", enforcer), handlers.GetEmployeeUsers())
		}

	}

	//------------
	//CUSTOMERS ROUTES
	//------------
	customersProtectedRoutes := apiRoutes.Group("/customers")
	{
		customersProtectedRoutes.GET("/", middleware.Authorize("rbac::data", "read", enforcer), handlers.GetAllCustomers())
		customersProtectedRoutes.POST("/", middleware.Authorize("rbac::data", "write", enforcer), handlers.AddCustomer())
		customersProtectedRoutes.DELETE("/:id", middleware.Authorize("rbac::data", "write", enforcer), handlers.DeleteCustomer(enforcer))
		customersProtectedRoutes.PUT("/", middleware.Authorize("rbac::data", "write", enforcer), handlers.UpdateCustomer())

		associations := customersProtectedRoutes.Group("/associations")
		{
			associations.GET("/:id", middleware.Authorize("rbac::data", "read", enforcer), handlers.GetCustomerUsers(enforcer))
			associations.PUT("/", middleware.Authorize("rbac::data", "read", enforcer), handlers.ToggleCustomerUserAccess(enforcer))
			associations.POST("/", middleware.Authorize("rbac::data", "write", enforcer), handlers.AddCustomerUserAssociation(enforcer))
			associations.DELETE("/", middleware.Authorize("rbac::data", "write", enforcer), handlers.DeleteCustomerUserAssociation(enforcer))
		}
	}

	//------------
	//FIREBASE ROUTES
	//------------
	firebaseProtectedRoutes := apiRoutes.Group("/firebase")
	{
		firebaseProtectedRoutes.GET("/", middleware.Authorize("rbac::data", "read", enforcer), handlers.GetAllFirebaseUsers())
		firebaseProtectedRoutes.POST("/", middleware.Authorize("rbac::data", "write", enforcer), handlers.AddFirebaseUser())
		firebaseProtectedRoutes.DELETE("/:id", middleware.Authorize("rbac::data", "write", enforcer), handlers.DeleteFirebaseUser(enforcer))
	}

	// SERVE FRONTEND
	if config.ENV("APP_ENV") == "prod" {
		fmt.Println("Production mode")
		log.Println("Production mode")
		gin.SetMode(gin.ReleaseMode)

		httpRouter.LoadHTMLGlob("../frontend/build/index.html")

		httpRouter.NoRoute(func(c *gin.Context) {
			httpRouter.Use(static.Serve("/", static.LocalFile("../frontend/build", true)))
			c.HTML(http.StatusOK, "index.html", nil)
		})
	}

	var envs map[string]string
	envs, err = godotenv.Read("./info/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := envs["APIPORT"]
	httpRouter.Run(":" + port)

}
