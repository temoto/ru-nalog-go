package umka

import (
	"fmt"
	"time"
)

type Status struct { //nolint:maligned
	AgentFlags     byte     `json:"agentFlags" fdn:"1057"`
	AllowGames     bool     `json:"allowGames" fdn:"1193"`
	AllowLotteries bool     `json:"allowLotteries" fdn:"1126"`
	AllowServices  bool     `json:"allowServices" fdn:"1109"`
	AtmNumber      string   `json:"atmNumber" fdn:"1036"`
	AutomatMode    bool     `json:"automatMode" fdn:"1001"`
	Cash           uint64   `json:"cash"`
	CashBoxNumber  uint32   `json:"cashBoxNumber" fdn:"1"`
	Cashier        uint32   `json:"cashier"`
	CycleNumber    uint32   `json:"cycleNumber" fdn:"1"`
	Dt             string   `json:"dt"`
	Email          string   `json:"email"`
	ExcisableGoods bool     `json:"excisableGoods" fdn:"1207"`
	ExternPrinter  bool     `json:"externPrinter" fdn:"1221"`
	FSFDFVersion   byte     `json:"fSFDFVersion" fdn:"1190"`
	FDFVersion     byte     `json:"fDFVersion"`
	Flags          byte     `json:"flags"`
	FnsSite        string   `json:"fnsSite" fdn:"1060"`
	FsNumber       string   `json:"fsNumber" fdn:"1041"`
	FsStatus       struct { //nolint:maligned
		CycleIsOpen   byte   `json:"cycleIsOpen"`
		DebugMode     bool   `json:"debugMode"`
		FsNumber      string `json:"fsNumber"`
		FsVersion     string `json:"fsVersion"`
		LastDocDt     string `json:"lastDocDt"`
		LastDocNumber uint64 `json:"lastDocNumber"`
		LifeTime      struct {
			AvailableRegistrations uint32 `json:"availableRegistrations"`
			CompletedRegistrations uint32 `json:"completedRegistrations"`
			ExpirationDt           string `json:"expirationDt"` // "2020-07-20"
		} `json:"lifeTime"`
		Phase     byte `json:"phase"`
		Transport struct {
			DocIsReading     bool   `json:"docIsReading"`
			FirstDocDt       string `json:"firstDocDt"`
			FirstDocNumber   uint64 `json:"firstDocNumber" fdn:"1116"`
			OfflineDocsCount uint32 `json:"offlineDocsCount" fdn:"1097"`
			State            uint32 `json:"state"`
		} `json:"transport"`
	} `json:"fsStatus"`
	InternetOnly     bool   `json:"internetOnly" fdn:"1108"`
	Introductions    uint32 `json:"introductions"`
	IntroductionsSum uint64 `json:"introductionsSum"`
	MakeBso          bool   `json:"makeBso"`
	Model            uint16 `json:"model"`
	Modelstr         string `json:"modelstr"` // "УМКА-01-ФА"
	OfdInn           string `json:"ofdInn" fdn:"1017"`
	OfdName          string `json:"ofdName" fdn:"1046"`
	OfflineMode      bool   `json:"offlineMode"`
	PaymentAddress   string `json:"paymentAddress" fdn:"1009"` // "г. Воронеж, ул. Липецкая, д.3"
	PaymentPlace     string `json:"paymentPlace" fdn:"1187"`   // "ОФИС1"
	Payouts          uint32 `json:"payouts"`
	PayoutsSum       uint64 `json:"payoutsSum"`
	RegCashierInn    string `json:"regCashierInn"`            // "000000000000"
	RegCashierName   string `json:"regCashierName"`           // "CASHIER 17"
	RegDate          string `json:"regDate"`                  // "2017-07-12"
	RegDocNumber     uint64 `json:"regDocNumber"`             // 1
	RegNumber        string `json:"regNumber" fdn:"1037"`     // "0000000001020321"
	ShortFlags       uint32 `json:"shortFlags"`               // 3
	Taxes            uint32 `json:"taxes" fdn:"1055"`         // 15
	UseEncryption    bool   `json:"useEncryption" fdn:"1056"` // false
	UserInn          string `json:"userInn" fdn:"1018"`       // "7725225244"
	UserName         string `json:"userName" fdn:"1048"`      // "ООО ВЕКТОР-М"
	Serial           string `json:"serial" fdn:"1013"`        // "16999987"

	Mode        string
	XXX_Mode    uint32 `json:"mode"`    // 0:choice 1:reg 2:x-report 3:z-report 4:prog 5:serial 6:fstore 7:aux
	XXX_SubMode uint32 `json:"subMode"` //

	// KktVersion string `fdn:"1188"`
	// KktFfdVersion string `fdn:"1189"`
	// FsFfdVersion string `fdn:"1190"`
	XXX_Ver    uint32 `json:"ver"`
	XXX_Subver uint32 `json:"subver"`
}

func (self *Status) IsCycleOpen() bool { return self.FsStatus.CycleIsOpen == 1 }

func (self *Status) FsExpireDate() time.Time {
	const layout = "2006-01-02"
	t, err := time.Parse(layout, self.FsStatus.LifeTime.ExpirationDt)
	if err != nil {
		panic(fmt.Sprintf("fs lifetime expire=%s err=%v", self.FsStatus.LifeTime.ExpirationDt, err))
	}
	return t
}

func (self *Status) OfdOfflineCount() uint32 {
	return self.FsStatus.Transport.OfflineDocsCount
}
