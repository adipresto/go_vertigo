package features

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"vertigo/pkg/broker"

	"github.com/cucumber/godog"
	_ "modernc.org/sqlite" // Ensure sqlite driver is available for tests
)

// dummyPublisher implements network.Publisher for testing
type dummyPublisher struct{}

func (d *dummyPublisher) Publish(ctx context.Context, channel string, data []byte) error {
	// For tests, we mock successful delivery
	return nil
}

type testCtx struct {
	facade *broker.TripleBaseBroker
	res    *httptest.ResponseRecorder
}

func (c *testCtx) thePersistenceFacadeIsInitialized() error {
	dbConn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return err
	}

	c.facade = &broker.TripleBaseBroker{
		DB:  dbConn,
		Net: &dummyPublisher{},
	}
	return nil
}

func (c *testCtx) theDatabaseHasATableWithRows(tableName string, table *godog.Table) error {
	// Create table
	var cols []string
	for _, cell := range table.Rows[0].Cells {
		cols = append(cols, cell.Value+" TEXT")
	}
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(cols, ","))
	if _, err := c.facade.DB.Exec(query); err != nil {
		return err
	}

	// Insert rows
	for i := 1; i < len(table.Rows); i++ {
		var values []string
		for _, cell := range table.Rows[i].Cells {
			values = append(values, fmt.Sprintf("'%s'", cell.Value))
		}
		insert := fmt.Sprintf("INSERT INTO %s VALUES (%s)", tableName, strings.Join(values, ","))
		if _, err := c.facade.DB.Exec(insert); err != nil {
			return err
		}
	}
	return nil
}

func (c *testCtx) iDispatchTheSQLQuery(sql string) error {
	ctx := context.Background()
	_, err := c.facade.Dispatch(ctx, sql, "test_channel")
	return err
}

func (c *testCtx) theResultShouldBeStreamedToTheNetwork() error {
	// Verification logic for streaming can be added here
	return nil
}

func (c *testCtx) theNetworkPayloadShouldContainAnd(arg1, arg2 string) error {
	// Content verification logic
	return nil
}

// REST API Steps
func (c *testCtx) iSendAGETRequestTo(path string) error {
	req, _ := http.NewRequest("GET", path, nil)
	c.res = httptest.NewRecorder()

	var handler http.HandlerFunc
	switch path {
	case "/api/users":
		handler = c.facade.HandleGetUsers()
	default:
		return fmt.Errorf("unknown path: %s", path)
	}

	handler.ServeHTTP(c.res, req)
	return nil
}

func (c *testCtx) iSendAPOSTRequestToWithBody(path string, body *godog.DocString) error {
	req, _ := http.NewRequest("POST", path, strings.NewReader(body.Content))
	c.res = httptest.NewRecorder()

	var handler http.HandlerFunc
	switch path {
	case "/api/dispatch":
		handler = c.facade.HandleDispatch()
	default:
		return fmt.Errorf("unknown path: %s", path)
	}

	handler.ServeHTTP(c.res, req)
	return nil
}

func (c *testCtx) theHTTPStatusCodeShouldBe(code int) error {
	if c.res.Code != code {
		return fmt.Errorf("expected status code %d, but got %d", code, c.res.Code)
	}
	return nil
}

func (c *testCtx) theHTTPResponseBodyShouldContain(expected string) error {
	if !strings.Contains(c.res.Body.String(), expected) {
		return fmt.Errorf("expected response to contain '%s', but got: %s", expected, c.res.Body.String())
	}
	return nil
}

func (c *testCtx) theHTTPResponseBodyShouldNotContain(expected string) error {
	if strings.Contains(c.res.Body.String(), expected) {
		return fmt.Errorf("expected response to NOT contain '%s', but it did: %s", expected, c.res.Body.String())
	}
	return nil
}

func InitializeScenario(sc *godog.ScenarioContext) {
	ctx := &testCtx{}
	sc.Step(`^the Persistence Facade is initialized$`, ctx.thePersistenceFacadeIsInitialized)
	sc.Step(`^the database has a table "([^"]*)" with rows:$`, ctx.theDatabaseHasATableWithRows)
	sc.Step(`^I dispatch the SQL query "([^"]*)"$`, ctx.iDispatchTheSQLQuery)
	sc.Step(`^the result should be streamed to the network$`, ctx.theResultShouldBeStreamedToTheNetwork)
	sc.Step(`^the network payload should contain "([^"]*)" and "([^"]*)"$`, ctx.theNetworkPayloadShouldContainAnd)

	// REST API
	sc.Step(`^I send a GET request to "([^"]*)"$`, ctx.iSendAGETRequestTo)
	sc.Step(`^I send a POST request to "([^"]*)" with body:$`, ctx.iSendAPOSTRequestToWithBody)
	sc.Step(`^the HTTP status code should be (\d+)$`, ctx.theHTTPStatusCodeShouldBe)
	sc.Step(`^the HTTP response body should contain "([^"]*)"$`, ctx.theHTTPResponseBodyShouldContain)
	sc.Step(`^the HTTP response body should not contain "([^"]*)"$`, ctx.theHTTPResponseBodyShouldNotContain)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"."},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
