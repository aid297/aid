package validatorV3

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/aid297/aid/ptr"
)

type (
	StringTest struct {
		Name1 string  `v-rule:"required;min>2;max<=10;not-in:张三" v-name:"姓名"`
		Name2 *string `v-rule:"required" v-name:"姓名1"`
		Name3 *string `v-rule:"required;min>0;in:王五,赵六"`
	}

	IntTest struct {
		Age1 int   `v-rule:"required"`
		Age2 *int  `v-rule:"required"`
		Age3 *int8 `v-rule:"required;min>=5"`
	}

	TimeTest struct {
		Time1 time.Time `json:"time1" v-rule:"datetime"`
	}
)

func Test1(t *testing.T) {
	st := StringTest{Name1: "张三", Name2: nil, Name3: ptr.New("")}
	t.Logf("%v", APP.Validator.Once().Checker(st).Validate().Wrongs())
}

func Test2(t *testing.T) {
	it := &IntTest{0, nil, ptr.New(int8(5))}

	t.Logf("%v", APP.Validator.Once().Checker(it).Validate().Wrongs())
}

func Test3(t *testing.T) {
	it := &IntTest{0, ptr.New(1), ptr.New(int8(5))}

	t.Logf("%v", APP.Validator.Once().Checker(it).Validate(func(data any) error {
		data.(*IntTest).Age2 = ptr.New(111)
		return nil
	}).Wrongs())

	t.Logf("%v", *it.Age2)
}

func Test4(t *testing.T) {
	tt := TimeTest{time.Now()}
	if err := json.Unmarshal([]byte(`{"time1": "2017-10-19T13:00:00+0200"}`), &tt); err != nil {
		t.Fatalf("反序列化失败： %v", err)
	}

	t.Logf("%v", APP.Validator.Once().Checker(tt).Validate().Wrongs())
}

type MainStoreRequest struct {
	Name           string  `json:"name" v-rule:"(required)(string)(min>1)(max<=30)" v-name:"项目名称"`
	Identification string  `json:"identification" v-rule:"(required)(string)(min>1)(max<=30)" v-name:"项目标识" v-ex:"PROJECT-IDENTITY"`
	Desc           *string `json:"desc" v-rule:"(string)(max<=100)" v-name:"项目描述"`
	TeamID         uint64  `json:"teamId" v-rule:"(required)(uint)(min>0)" v-name:"团队id"`
	OwnerUUID      string  `json:"ownerUuid" v-rule:"(required)(string)(size=36)" v-name:"项目所有者uuid"`
	OwnerUsername  string  `json:"ownerUsername" v-rule:"(required)(string)(min>1)(max<=64)" v-name:"项目所有者昵称"`
}

func Test5(t *testing.T) {
	j := `{
    "name": "项目1",
    "identification": "project-1",
    "desc": null,
    "teamId": 2,
    "ownerUuid": "027c12c1-7524-4439-9267-c351b6d4a9aa",
    "ownerNickname": "test-user-a"
}`
	m := MainStoreRequest{}
	if err := json.Unmarshal([]byte(j), &m); err != nil {
		t.Errorf("反序列化失败： %v", err)
		return
	}

	t.Logf("%v", APP.Validator.Once().Checker(m).Validate().Wrongs())
}
