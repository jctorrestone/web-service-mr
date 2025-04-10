package model

type Data any

type Response struct {
	Data     Data  `json:"data"`
	PrevPage int   `json:"prev_page"`
	Page     int   `json:"page"`
	NextPage int   `json:"next_page"`
	LastPage int   `json:"last_page"`
	Total    int64 `json:"total"`
}

type FullRecord struct {
	RecordObj       Record            `json:"record"`
	DiseasesHistory []DiseaseHistory  `json:"diseases_history"`
	Symptoms        []Symptom         `json:"symptoms"`
	VitalSigns      []RecordVitalSign `json:"vital_signs"`
	Diseases        []Disease         `json:"idx"`
	Exams           []Exam            `json:"exams"`
	Treatments      []Treatment       `json:"treatments"`
}

type Patient struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Lastname string `json:"last_name"`
	Gender   bool   `json:"gender"`
}

type Exam struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
}

type Record struct {
	ID         int64   `json:"id"`
	Category   string  `json:"category"`
	PrimaryID  int64   `json:"primary_record_id"`
	PatientObj Patient `json:"patient"`
	Date       string  `json:"rdate"`
	Age        int64   `json:"age"`
	Weight     int64   `json:"weight"`
	Height     int64   `json:"height"`
	Duration   int64   `json:"duration"`
}

type RecordExam struct {
	RecordID  int64 `json:"record_id"`
	RecordObj Record
	ExamID    int64 `json:"exam_id"`
	ExamObj   Exam
}

type Unit struct {
	ID          int64  `json:"id"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
}

type VitalSign struct {
	ID          int64  `json:"id"`
	UnitObj     Unit   `json:"unit"`
	Description string `json:"description"`
}

type RecordVitalSign struct {
	RecordID    int64   `json:"record_id"`
	VitalSignID int64   `json:"vital_sign_id"`
	Description string  `json:"vital_sign_desc"`
	UnitID      int64   `json:"unit_id"`
	Symbol      string  `json:"unit_symbol"`
	Value       float64 `json:"value"`
}

type Disease struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
}

type DiseaseHistory struct {
	RecordID    int64  `json:"record_id"`
	DiseaseID   int64  `json:"disease_id"`
	DiseaseDesc string `json:"disease_desc"`
	Description string `json:"description"`
}

type Symptom struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
}

type RecordSymptom struct {
	RecordID   int64 `json:"record_id"`
	RecordObj  Record
	SymptomID  int64 `json:"symptom_id"`
	SymptomObj Symptom
}

type Idx struct {
	RecordID   int64 `json:"record_id"`
	RecordObj  Record
	DiseaseID  int64 `json:"disease_id"`
	DiseaseObj Symptom
}

type Shape struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
}

type Formulation struct {
	ID       int64 `json:"id"`
	ShapeObj Shape `json:"shape"`
	UnitObj  Unit  `json:"unit"`
}

type Medicine struct {
	ID             int64       `json:"id"`
	FormulationObj Formulation `json:"formulation"`
	Name           string      `json:"name"`
	Dose           int64       `json:"dose"`
}

type Treatment struct {
	RecordID      int64   `json:"record_id"`
	MedicineID    int64   `json:"medicine_id"`
	Name          string  `json:"medicine_name"`
	Dose          int64   `json:"medicine_dose"`
	FormulationID int64   `json:"formulation_id"`
	ShapeID       int64   `json:"shape_id"`
	Description   string  `json:"shape_description"`
	UnitID        int64   `json:"unit_id"`
	Symbol        string  `json:"unit_symbol"`
	Quantity      int64   `json:"quantity"`
	Dosage        float64 `json:"dosage"`
	Frequency     int64   `json:"frequency"`
	Instructions  string  `json:"instructions"`
}
