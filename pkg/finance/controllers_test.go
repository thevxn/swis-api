package finance

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.savla.dev/swis/v5/pkg/core"

	"github.com/stretchr/testify/assert"
)

// app = array of pointers to pointers to Cache
type appCache []**core.Cache

var TestPackage *core.Package = &core.Package{
	Name: pkgName,
	Cache: []**core.Cache{
		&CacheAccounts,
		&CacheItems,
	},
	Routes: Routes,
}

/*
 *  unit/integration tests
 */

/*
 *  common
 */

func TestPostDumpRestore(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var items = struct {
		Accounts map[string]Account `json:"accounts"`
		Items    map[string]Item    `json:"items"`
	}{
		Accounts: map[string]Account{
			"test_acc": {
				ID:            "test_acc",
				AccountNumber: "660195318",
				Currency:      "EUR",
				SWIFT:         "CZxxx",
				IBAN:          "999",
				Owner:         "unknown",
			},
			"": {
				ID:            "",
				AccountNumber: "006978132",
			},
			"invalid": {
				ID:            "",
				AccountNumber: "066566998",
			},
		},
		Items: map[string]Item{
			"test_item": {
				ID:          "test_item",
				Type:        "expense",
				Amount:      55.66,
				Currency:    "EUR",
				AccountID:   "test_acc",
				PaymentDate: time.Now(),
			},
		},
	}

	jsonValue, _ := json.Marshal(items)
	req, _ := http.NewRequest("POST", "/finance/restore", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Count []int `json:"counter"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusCreated, w.Code)

	// this is not good, number 2 here means, that the 'invalid' account is imported too,
	// even though it does not suffice given property requirements...
	// TODO: implement simple JSON binding, or some sort of an item verificator
	assert.Equal(t, []int{2, 1}, ret.Count)
}

func TestGetRootData(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/finance/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var items = struct {
		Accounts map[string]Account `json:"accounts"`
		Items    map[string]Item    `json:"items"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &items)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, items.Accounts)
	assert.NotEmpty(t, items.Items)
}

/*
 *  accounts
 */

func TestPostNewAccountByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var acc Account = Account{
		ID:            "test_acc",
		AccountNumber: "660195318",
		Currency:      "EUR",
		SWIFT:         "CZxxx",
		IBAN:          "999",
		Owner:         "unknown",
	}

	jsonValue, _ := json.Marshal(acc)
	req, _ := http.NewRequest("POST", "/finance/accounts/test_acc", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Account Account `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	// this account has already been imported/created using the dump restore test
	// --- thus resulting in HTTP/409 Conflict
	//assert.NotEqual(t, "", item.Account.AccountNumber)
	//assert.NotEqual(t, "", item.Account.ID)
	//assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestGetAccounts(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/finance/accounts/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var items = struct {
		Accounts map[string]Account `json:"items"`
		Count    int                `json:"count"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &items)

	assert.Equal(t, http.StatusOK, w.Code)

	// number 2 down below is wrong, see notes above (cca L87)
	assert.Equal(t, 2, items.Count)
	assert.NotEmpty(t, items.Accounts)
}

func TestGetAccountByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/finance/accounts/test_acc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Account Account `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, item.Account)
	assert.NotEqual(t, "", item.Account.AccountNumber)
}

func TestGetAccountByOwnerKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/finance/accounts/owner/fantomas", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var items = struct {
		Accounts map[string]Account `json:"items"`
		Count    int                `json:"count"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &items)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 0, items.Count)
	assert.Empty(t, items.Accounts)
}

func TestUpdateAccountByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var acc Account = Account{
		ID:            "test_acc",
		AccountNumber: "660195318",
		Currency:      "EUR",
		SWIFT:         "CZxxx",
		IBAN:          "999",
		Owner:         "fantomas",
	}

	jsonValue, _ := json.Marshal(acc)
	req, _ := http.NewRequest("PUT", "/finance/accounts/test_acc", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var item = struct {
		Account Account `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &item)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, "", item.Account.AccountNumber)
	assert.Equal(t, "fantomas", item.Account.Owner)
}

func TestDeleteAccountByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("DELETE", "/finance/accounts/test_acc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Key string `json:"key"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test_acc", ret.Key)
}

/*
 *  items (financial)
 */

func TestPostNewItemByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var item Item = Item{
		ID:          "test_item2",
		Type:        "income",
		Amount:      5.77,
		Currency:    "EUR",
		AccountID:   "test_acc",
		PaymentDate: time.Now(),
	}

	jsonValue, _ := json.Marshal(item)
	req, _ := http.NewRequest("POST", "/finance/items/test_item2", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetItems(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/finance/items/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var items = struct {
		Items map[string]Item `json:"items"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &items)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, items.Items)
}

func TestGetItemsByAccountID(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("GET", "/finance/items/account/test_acc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var items = struct {
		Items map[string]Item `json:"items"`
		Count int             `json:"count"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &items)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 2, items.Count)
	assert.NotEmpty(t, items.Items)
}

func TestUpdateItemByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	var item Item = Item{
		ID:          "test_item2",
		Type:        "income",
		Amount:      6.88,
		Currency:    "EUR",
		AccountID:   "test_acc",
		PaymentDate: time.Now(),
	}

	jsonValue, _ := json.Marshal(item)
	req, _ := http.NewRequest("PUT", "/finance/items/test_item2", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Item Item `json:"item"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 6.88, ret.Item.Amount)
}

func TestDeleteItemByKey(t *testing.T) {
	r := core.SetupTestEnv(TestPackage)

	req, _ := http.NewRequest("DELETE", "/finance/items/test_item2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var ret = struct {
		Key string `json:"key"`
	}{}
	json.Unmarshal(w.Body.Bytes(), &ret)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test_item2", ret.Key)
}
