package model

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
	ID         int64 `json:"id"`
	PatientID  int64 `json:"patient_id"`
	PatientObj Patient
	Date       string `json:"rdate"`
	Age        int64  `json:"age"`
	Weight     int64  `json:"weight"`
	Height     int64  `json:"height"`
	Duration   int64  `json:"duration"`
}

type RecordExam struct {
	ID        int64 `json:"id"`
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
	ID          int64 `json:"id"`
	UnitID      int64 `json:"unit_id"`
	UnitObj     Unit
	Description string `json:"description"`
}

type RecordVitalSign struct {
	ID           int64 `json:"id"`
	RecordID     int64 `json:"record_id"`
	RecordObj    Record
	VitalSignID  int64 `json:"vital_sign_id"`
	VitalSignObj VitalSign
	Value        float64 `json:"value"`
}

type Disease struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
}

type DiseaseHistory struct {
	ID          int64 `json:"id"`
	RecordID    int64 `json:"record_id"`
	RecordObj   Record
	DiseaseID   int64 `json:"disease_id"`
	DiseaseObj  Disease
	Description string `json:"description"`
}

type Symptom struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
}

type RecordSymptom struct {
	ID         int64 `json:"id"`
	RecordID   int64 `json:"record_id"`
	RecordObj  Record
	SymptomID  int64 `json:"symptom_id"`
	SymptomObj Symptom
}

type Idx struct {
	ID         int64 `json:"id"`
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
	ShapeID  int64 `json:"shape_id"`
	ShapeObj Shape
	UnitID   int64 `json:"unit_id"`
	UnitObj  Unit
}

type Medicine struct {
	ID             int64 `json:"id"`
	FormulationID  int64 `json:"formulation_id"`
	FormulationObj Formulation
	Name           string `json:"name"`
	Dose           int64  `json:"dose"`
}

type Treatment struct {
	ID           int64 `json:"id"`
	RecordID     int64 `json:"record_id"`
	RecordObj    Record
	MedicineID   int64 `json:"medicine_id"`
	MedicineObj  Medicine
	Quantity     int64   `json:"quantity"`
	Dosage       float64 `json:"dosage"`
	Frequency    int64   `json:"frequency"`
	Instructions string  `json:"instructions"`
}
