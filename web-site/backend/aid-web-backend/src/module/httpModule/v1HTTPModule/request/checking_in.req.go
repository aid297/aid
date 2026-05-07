package request

type CheckingInCalRequest struct {
	OriginalDate      string   `form:"originalDate" v-role:"(required)(not-empty)" v-name:"起始时间"`
	ForceWorkdays     []string `form:"forceWorkdays"`
	ForceAnnualLeaves []string `form:"forceAnnualLeaves"`
	Holidays2         []string `form:"holidays2"`
	Holidays3         []string `form:"holidays3"`
}
