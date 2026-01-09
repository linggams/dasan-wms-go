package domain

import "time"

// Stage represents the fabric processing stage
type Stage string

const (
	StageInventory      Stage = "inventory"
	StageRelaxation     Stage = "relaxation"
	StageCuttingWIP     Stage = "cutting_wip"
	StageStockFabric    Stage = "stock_fabric"
	StageCNCM           Stage = "cncm"
	StageWashing        Stage = "washing"
	StageReturnSupplier Stage = "return_supplier"
	StageDestroy        Stage = "destroy"
	StageQCFabric       Stage = "qc_fabric"
)

// StageInfo represents stage information for overview
type StageInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetAllStages returns all available stages
func GetAllStages() []StageInfo {
	return []StageInfo{
		{ID: 1, Name: string(StageInventory)},
		{ID: 2, Name: string(StageCuttingWIP)},
		{ID: 3, Name: string(StageStockFabric)},
		{ID: 4, Name: string(StageCNCM)},
		{ID: 5, Name: string(StageWashing)},
		{ID: 6, Name: string(StageReturnSupplier)},
		{ID: 7, Name: string(StageDestroy)},
		{ID: 8, Name: string(StageRelaxation)},
		{ID: 9, Name: string(StageQCFabric)},
	}
}

// IsValidStage checks if the stage is valid
func IsValidStage(s string) bool {
	validStages := map[string]bool{
		string(StageInventory):      true,
		string(StageRelaxation):     true,
		string(StageCuttingWIP):     true,
		string(StageStockFabric):    true,
		string(StageCNCM):           true,
		string(StageWashing):        true,
		string(StageReturnSupplier): true,
		string(StageDestroy):        true,
		string(StageQCFabric):       true,
	}
	return validStages[s]
}

// User represents the user entity
type User struct {
	ID                   int64      `json:"id"`
	Name                 string     `json:"name"`
	Email                string     `json:"email"`
	EmailVerifiedAt      *time.Time `json:"email_verified_at,omitempty"`
	Password             string     `json:"-"` // Don't serialize password
	TwoFactorSecret      *string    `json:"two_factor_secret,omitempty"`
	TwoFactorRecovery    *string    `json:"two_factor_recovery_codes,omitempty"`
	TwoFactorConfirmedAt *time.Time `json:"two_factor_confirmed_at,omitempty"`
	RememberToken        *string    `json:"remember_token,omitempty"`
	LastMobileToken      *string    `json:"last_mobile_token,omitempty"`
	PasswordChangedAt    *time.Time `json:"password_changed_at,omitempty"`
	IsActive             int        `json:"is_active"`
	EmailForgotPassword  *string    `json:"email_forgot_password,omitempty"`
	ProfilePhotoPath     *string    `json:"profile_photo_path,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at,omitempty"`
}

// Fabric represents the fabric entity
type Fabric struct {
	ID                 int64      `json:"id"`
	Code               string     `json:"code"`
	FabricIncomingID   *int64     `json:"fabric_incoming_id,omitempty"`
	SupplierID         *int64     `json:"supplier_id,omitempty"`
	Color              string     `json:"color,omitempty"`
	Lot                string     `json:"lot,omitempty"`
	Roll               string     `json:"roll,omitempty"`
	Weight             string     `json:"weight,omitempty"`
	Width              *string    `json:"width,omitempty"`
	Yard               string     `json:"yard"`
	UnitID             *int64     `json:"unit_id,omitempty"`
	FabricType         *string    `json:"fabric_type,omitempty"`
	FabricContain      *string    `json:"fabric_contain,omitempty"`
	RackID             *int64     `json:"rack_id,omitempty"`
	BlockID            *int64     `json:"block_id,omitempty"`
	RelaxationRackID   *int64     `json:"relaxation_rack_id,omitempty"`
	RelaxationBlockID  *int64     `json:"relaxation_block_id,omitempty"`
	FinishDate         *string    `json:"finish_date,omitempty"`
	QCResult           *string    `json:"qc_result,omitempty"`
	Status             *string    `json:"status,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
	Buyer              string     `json:"buyer"`
	Style              string     `json:"style"`
	YardStock          string     `json:"yard_stock,omitempty"`
	YardRemaining      string     `json:"yard_remaining,omitempty"`
	FabricIncoming     *FabricIncoming `json:"fabric_incoming,omitempty"`
	Block              *Block     `json:"block,omitempty"`
	Rack               *Rack      `json:"rack,omitempty"`
	Inventory          *Inventory `json:"inventory,omitempty"`
}

// FabricIncoming represents incoming fabric shipment
type FabricIncoming struct {
	ID          int64     `json:"id"`
	Code        string    `json:"code"`
	Datetime    string    `json:"datetime"`
	BuyerID     int64     `json:"buyer_id"`
	Style       string    `json:"style"`
	VendorID    int64     `json:"vendor_id"`
	StickerPath *string   `json:"sticker_path,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	TotalYards  float64   `json:"total_yards"`
	TotalItems  int       `json:"total_items"`
	Buyer       *Buyer    `json:"buyer,omitempty"`
}

// Buyer represents buyer entity
type Buyer struct {
	ID            int64      `json:"id"`
	Code          string     `json:"code"`
	Name          string     `json:"name"`
	ContactPerson *string    `json:"contact_person,omitempty"`
	Phone         *string    `json:"phone,omitempty"`
	Address       *string    `json:"address,omitempty"`
	Email         *string    `json:"email,omitempty"`
	City          *string    `json:"city,omitempty"`
	Country       *string    `json:"country,omitempty"`
	Status        int        `json:"status"`
	Remarks       *string    `json:"remarks,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

// Inventory represents fabric inventory tracking
type Inventory struct {
	ID        int64      `json:"id"`
	Datetime  string     `json:"datetime"`
	RefNumber *string    `json:"ref_number,omitempty"`
	FabricID  int64      `json:"fabric_id"`
	Stage     string     `json:"stage"`
	Remarks   *string    `json:"remarks,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// Block represents storage block
type Block struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// Rack represents storage rack
type Rack struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
