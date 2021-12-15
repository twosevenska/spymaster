package mongo

import (
	"errors"
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Config holds the required configuration for a MongoDB connection
type Config struct {
	Hosts    []string `envconfig:"hosts_tls" default:"0.0.0.0:27017"`
	Database string   `envconfig:"-" default:"spymaster"`
	User     string   `envconfig:"-" default:"spymaster"`
	Password string   `envconfig:"-" default:"face1t"`
}

// Client represents a MongoDB client
type Client struct {
	Database *mgo.Database
	session  *mgo.Session
	db       string
}

var (
	// ErrInvalidUUID indicates that an invalid UUID is passed as a parameter
	ErrInvalidUUID = errors.New("invalid database ID")
)

const (
	defaultMaxQueryTime = 2 * time.Second
)

// Connect connects to a MongoDB cluster and returns a client
func Connect(conf Config) (*Client, error) {
	log.Printf("Connecting to MongoDB @ %s", conf.Hosts)

	dialInfo := &mgo.DialInfo{
		Addrs:    conf.Hosts,
		Database: conf.Database,
		Username: conf.User,
		Password: conf.Password,
		Timeout:  time.Second * 10,
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s", err)
		return nil, err
	}
	ensureIndices(session, conf.Database)
	session.SetMode(mgo.Monotonic, true)
	database := session.DB(conf.Database)
	return &Client{
		Database: database,
		session:  session,
		db:       conf.Database,
	}, nil
}

func ensureIndices(s *mgo.Session, db string) {
	session := s.Copy()
	defer session.Close()
	c := session.DB(db).C("rules")
	createIndex(c, []string{"_id", "country"}, false)
	createIndex(c, []string{"nickname", "email"}, true)
	createIndex(c, []string{"nickname", "country"}, false)
}

func createIndex(c *mgo.Collection, keys []string, unique bool) {
	i := mgo.Index{
		Key:        keys,
		Unique:     unique,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(i)
	if err != nil {
		log.Fatal(err)
	}
}

// Copy creates a new session by calling session.Copy on the initial session obtained at dial time
func (c Client) Copy() Client {
	session := c.session.Copy()
	database := session.DB(c.db)
	return Client{
		Database: database,
		session:  session,
		db:       c.db,
	}
}

// Close closes the MongoDB session
func (c Client) Close() {
	c.session.Close()
}

// IsDup returns whether err informs of a duplicate key error because a primary key index
// or a secondary unique index already has an entry with the given value.
func (c Client) IsDup(err error) bool {
	return mgo.IsDup(err)
}

// IsNotFound returns whether err informs of documents not matching the search criteria
func (c Client) IsNotFound(err error) bool {
	return err == mgo.ErrNotFound
}

// IsInvalidUUID returns whether err informs of an invalid UUID
func (c Client) IsInvalidUUID(err error) bool {
	return err == ErrInvalidUUID
}

// safeFind is a wrapper around the regular mgo driver find function that adds a query execution timeout
func safeFind(cl *mgo.Collection, query bson.M) *mgo.Query {
	return cl.Find(query).SetMaxTime(defaultMaxQueryTime)
}

// Ping runs a trivial ping command just to get in touch with the server.
func (c Client) Ping() error {
	return c.session.Ping()
}
