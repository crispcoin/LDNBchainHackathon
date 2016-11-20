package hackathon

import(
	"fmt"
	"error"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"strconv"
)

type dentist struct {
	Id string `json:"id"`                        //
	Given_name string `json:"given_name"`
	Surname string `json:"surname"`
	Gender string `json:"gender"`
	Status string `json:"status"`
	Year_of_qualification string `json:"Year_of_qual"`
	Diagnoses []diagnosis `json:"Diagnoses"`
}

func (den *dentist) AddDiagnose(d diagnosis) []diagnosis {
	den.Diagnoses = append(den.Diagnoses, d)
	return den.Diagnoses;
}


//TODO: ISSUE with moving from one slice to another !!!!
type patient struct {
	Id string `json:"id"`                        //NHS number in UK
	Name string `json:"name"`
	Address string `json:"address"`
	Open_diagnoses []diagnosis `json:"Open"`       //list of open cases
	Closed_diagnoses []diagnosis `json:"Close"`    //list of closed cases
}
func (p *patient) AddOpenDiagnose(d diagnosis) []diagnosis {
	p.Open_diagnoses = append(p.Open_diagnoses, d)
	return p.Open_diagnoses;
}
func (p *patient) AddClosedDiagnose(d diagnosis) []diagnosis {
	p.Closed_diagnoses = append(p.Closed_diagnoses, d)
	return p.Closed_diagnoses;
}

type diagnosis struct{
	Id string `json:"id"`
	//maybe unnecessary duplication
	patient_ID string `json:"patient_ID"`        //patient the treatment for
	dentist_ID string `json:"dentist_ID"`        //dentist suggesting the treatment
	Date string `json:"date"`
	ICD string `json:"icd"`                      //international classification diseases
 	drug_treatment string `json:"drug"`          //suggested treatment
 	data string `json:"data"`                    //anything to support the dentist's decision
	approved_by []dentist `json:"approved_by"`   //list of dentists approving the suggested treatment
	disapproved_by []dentist `json:"disapproved_by"`                     //list of dentists NON-approving the suggested treatment
}
func (d *diagnosis) AddApprovedBy(den dentist) []dentist {
	d.approved_by = append(d.approved_by, den)
	return d.approved_by;
}
func (d *diagnosis) AddDisapprovedBy(den dentist) []dentist {
	d.disapproved_by = append(d.disapproved_by, den)
	return d.disapproved_by;
}

type MedicalChaincode struct{
}

func (t *MedicalChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	// empty init ?

  	return nil,nil
}

//functions to add a entities
func (t *MedicalChaincode) addDentist(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	if len(args) %6 != 0  {
		return nil, errors.New("Incorrect number of args. Needs to be multiples of 6: (Id,Given_name,Surname,Gender,Status,Year_of_qualification")
	}

	DentistOBJ := dentist{}

	for index := 0; index < len(args); index += 6 {
		id = args[index];
		//TODO: look into this duplication, any advantage?
		DentistOBJ.Id = args[index];
		DentistOBJ.Given_name = args[index + 1];
		DentistOBJ.Surname = args[index + 2];
		DentistOBJ.Gender = args[index + 3];
		DentistOBJ.Status = args[index + 4];
		DentistOBJ.Year_of_qualification = args[index + 5];
		DentistOBJ.Diagnoses = []diagnosis{};

		DentistJSON, err := json.Marshal(DentistOBJ);
		if err != nil || DentistJSON == nil {
			return nil, errors.New("Converting entity struct to DentistJSON failed")
		}

		err = stub.PutState(id , DentistJSON)
		if err != nil {
			fmt.Printf("Error: could not update ledger")
			return nil, err
		}
	}
	return nil,nil
}

func (t *MedicalChaincode) addPatient(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	if len(args) %3 != 0  {
		return nil, errors.New("Incorrect number of args. Needs to be multiples of 3: (Id,name,address")
	}

	PatientOBJ := patient{}

	for index := 0; index < len(args); index += 3 {
		id = args[index];
		//TODO: look into this duplication, any advantage?
		PatientOBJ.Id = args[index];
		PatientOBJ.Name = args[index + 1];
		PatientOBJ.Address = args[index + 2];
		PatientOBJ.Open_diagnoses = []diagnosis{};
		PatientOBJ.Closed_diagnoses = []diagnosis{};

		PatientJSON, err := json.Marshal(PatientOBJ);
		if err != nil || PatientJSON == nil {
			return nil, errors.New("Converting entity struct to PatientJSON failed")
		}

		err = stub.PutState(id , PatientJSON)
		if err != nil {
			fmt.Printf("Error: could not update ledger")
			return nil, err
		}
	}
	return nil,nil
}


//TODO: store the diag in patient and dentist tooooooooooooo
func (t *MedicalChaincode) addDiagnosis(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) %7 != 0  {
		return nil, errors.New("Incorrect number of args. Needs to be multiples of 7: (Id,patient_ID,dentist_ID,date,ICD,treatment,data")
	}

	DiagnosisOBJ := diagnosis{}

	for index := 0; index < len(args); index += 6 {
		id = args[index];
		//TODO: look into this duplication, any advantage?
		DiagnosisOBJ.Id = args[index];
		patient_ID = args[index + 1];
		DiagnosisOBJ.patient_ID = patient_ID;
		dentist_ID = args[index + 2];
		DiagnosisOBJ.dentist_ID = dentist_ID;
		DiagnosisOBJ.Date = args[index + 3];
		DiagnosisOBJ.ICD = args[index + 4];
		DiagnosisOBJ.drug_treatment = args[index + 5];
		DiagnosisOBJ.data = args[index + 6];
		DiagnosisOBJ.approved_by = []dentist{};
		DiagnosisOBJ.disapproved_by = []dentist{};

		DiagnosisJSON, err := json.Marshal(DiagnosisOBJ);
		if err != nil || DiagnosisJSON == nil {
			return nil, errors.New("Converting entity struct to PatientJSON failed")
		}

		err = stub.PutState(id , DiagnosisJSON)
		if err != nil {
			fmt.Printf("Error: could not update ledger")
			return nil, err
		}

		//get the patient by ID
		patientJSON, err := stub.GetState(patient_ID)
		if patientJSON == nil {
			return nil, errors.New("Error: No account exists for user.")
		}

		patientOBJ := patient{}
		err = json.Unmarshal(patientJSON, &patientOBJ)
		if err != nil {
			return nil, errors.New("Invalid entity data pulled from ledger.")
		}

		//save the diagnose ID to his list
		patientOBJ.AddOpenDiagnose(id);
		newPatientJSON, err := json.Marshal(patientOBJ);
		if err != nil || newPatientJSON == nil {
			return nil, errors.New("Converting entity struct to PatientJSON failed")
		}

		//write it back to ledger
		err = stub.PutState(patient_ID , newPatientJSON)
		if err != nil {
			fmt.Printf("Error: could not update ledger")
			return nil, err
		}



		//get the dentist by ID
		dentistJSON, err := stub.GetState(dentist_ID)
		if dentistJSON == nil {
			return nil, errors.New("Error: No account exists for user.")
		}

		dentistOBJ := dentist{}
		err = json.Unmarshal(dentistJSON, &dentistOBJ)
		if err != nil {
			return nil, errors.New("Invalid entity data pulled from ledger.")
		}

		//save the diagnose ID to his list
		dentistOBJ.AddDiagnose(id)
		newDentistJSON, err := json.Marshal(dentistOBJ);
		if err != nil || newDentistJSON == nil {
			return nil, errors.New("Converting entity struct to PatientJSON failed")
		}

		//write it back to ledger
		err = stub.PutState(dentist_ID , newDentistJSON)
		if err != nil {
			fmt.Printf("Error: could not update ledger")
			return nil, err
		}
	}
	return nil,nil
}

func (t *MedicalChaincode) approveDiagnosis(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return nil,nil
}






//functions to delete entities

//invoke to call the above functions
func (t *MedicalChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return nil,nil
}


func (t *MedicalChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return nil,nil
}


func main() {
	err := shim.Start(new(MedicalChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}