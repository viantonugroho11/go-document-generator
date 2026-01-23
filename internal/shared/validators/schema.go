package validators

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
)

var (
	cacheOnce   sync.Once
	schemaCache *lru.Cache[string, *jsonschema.Schema]
	schemaCacheMu sync.RWMutex
)

func getSchemaCache() *lru.Cache[string, *jsonschema.Schema] {
	cacheOnce.Do(func() {
		// kapasitas dapat disesuaikan
		c, _ := lru.New[string, *jsonschema.Schema](128)
		schemaCache = c
	})
	return schemaCache
}

func ValidateSchema(schema map[string]any, payload map[string]any) error {
	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return err
	}
	keyBytes := sha256.Sum256(schemaBytes)
	key := hex.EncodeToString(keyBytes[:])

	cache := getSchemaCache()
	var sch *jsonschema.Schema
	// fast path with RLock
	schemaCacheMu.RLock()
	if v, ok := cache.Get(key); ok {
		schemaCacheMu.RUnlock()
		sch = v
	} else {
		schemaCacheMu.RUnlock()
		// compile dan cache
		c := jsonschema.NewCompiler()
		uri := "inmem://" + key + ".json"
		if err := c.AddResource(uri, bytes.NewReader(schemaBytes)); err != nil {
			return err
		}
		compiled, err := c.Compile(uri)
		if err != nil {
			return err
		}
		// write path
		schemaCacheMu.Lock()
		cache.Add(key, compiled)
		schemaCacheMu.Unlock()
		sch = compiled
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return sch.Validate(bytes.NewReader(payloadBytes))
}
