package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"vertigo/pkg/broker"
	"vertigo/pkg/config"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	_ "modernc.org/sqlite"
)

var globalBroker *broker.TripleBaseBroker

func main() {
	fmt.Println("🚀 Starting Vertigo (REST & GraphQL Demo)...")

	// 1. Load Configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Fatal: Failed to load config.yaml: %v", err)
	}

	// 2. Initialize Facade
	globalBroker, err = broker.NewBroker(cfg)
	if err != nil {
		log.Fatalf("Fatal: Failed to initialize Vertigo: %v", err)
	}

	// 3. Setup Schema for Demo
	setupSchema(globalBroker)

	// 4. Setup GraphQL
	gqlSchema, err := setupGraphQL()
	if err != nil {
		log.Fatal(err)
	}

	// 4. Register Routes
	http.HandleFunc("/", landingPage)
	http.HandleFunc("/api/users", globalBroker.HandleGetUsers())
	http.HandleFunc("/api/dispatch", globalBroker.HandleDispatch())
	// GraphQL Handler
	h := handler.New(&handler.Config{
		Schema:   &gqlSchema,
		Pretty:   true,
		GraphiQL: true,
	})
	http.Handle("/graphql", h)

	fmt.Printf("🌍 Server running at http://localhost:%d\n", cfg.Server.Port)
	fmt.Printf("📊 GraphiQL available at http://localhost:%d/graphql\n", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), nil))
}

func setupSchema(b *broker.TripleBaseBroker) {
	_, err := b.DB.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, email TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	// Insert dummy data if empty
	var count int
	b.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count == 0 {
		_, _ = b.DB.Exec("INSERT INTO users (name, email) VALUES (?, ?), (?, ?)",
			"Alice", "alice@example.com", "Bob", "bob@example.com")
	}
}

// GraphQL Setup
func setupGraphQL() (graphql.Schema, error) {
	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":    &graphql.Field{Type: graphql.Int},
			"name":  &graphql.Field{Type: graphql.String},
			"email": &graphql.Field{Type: graphql.String},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type: graphql.NewList(userType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					data, err := globalBroker.Dispatch(p.Context, "SELECT * FROM users", "graphql_users")
					if err != nil {
						return nil, err
					}
					
					// Decode JSON payload from phgo
					var payload struct {
						Data string `json:"data"`
					}
					if err := json.Unmarshal(data, &payload); err != nil {
						return nil, err
					}

					// Now unmarshal the nested data string
					var users []map[string]interface{}
					if err := json.Unmarshal([]byte(payload.Data), &users); err != nil {
						return nil, err
					}
					return users, nil
				},
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"dispatchSQL": &graphql.Field{
				Type: graphql.String, // Return raw JSON string
				Args: graphql.FieldConfigArgument{
					"sql":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"channel": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					sql := p.Args["sql"].(string)
					channel := p.Args["channel"].(string)
					data, err := globalBroker.Dispatch(p.Context, sql, channel)
					if err != nil {
						return nil, err
					}
					return string(data), nil
				},
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
}

func landingPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, `
		<html>
			<head>
				<title>phgo Demo</title>
				<style>
					body { font-family: sans-serif; padding: 20px; line-height: 1.6; }
					code { background: #f4f4f4; padding: 2px 5px; border-radius: 3px; }
					.container { max-width: 800px; margin: 0 auto; }
				</style>
			</head>
			<body>
				<div class="container">
					<h1>📦 phgo: Persistence Handler Demo</h1>
					<p>This is a realistic demo of the <strong>Persistence Handler Go</strong> package.</p>
					
					<h2>🚀 API Endpoints</h2>
					<ul>
						<li><code>GET /api/users</code> - Fetch all users via phgo</li>
						<li><code>POST /api/dispatch</code> - Generic SQL dispatcher</li>
						<li><code>GET /graphql</code> - GraphiQL Interactive Interface</li>
					</ul>

					<h2>📊 Sample GraphQL Query</h2>
					<pre><code>
query {
  users {
    id
    name
    email
  }
}
					</code></pre>

					<h2>📡 Sample REST Dispatch</h2>
					<pre><code>
curl -X POST http://localhost:8080/api/dispatch \
  -H "Content-Type: application/json" \
  -d '{"sql": "INSERT INTO users (name, email) VALUES (\"Charlie\", \"charlie@example.com\")", "channel": "user_updates"}'
					</code></pre>
				</div>
			</body>
		</html>
	`)
}
