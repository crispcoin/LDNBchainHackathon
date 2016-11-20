package hackathon

import(
	"fmt"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)

type dentist struct {
	Id string `json:"id"`                        //
	Given_name string `json:"given_name"`
	Surname string `json:"surname"`
	Gender string `json:"gender"`
	Status string `json:"status"`
	Year_of_qualification string `json:"Year_of_qual"`
	Diagnoses []string `json:"Diagnoses"`
}

func (den *dentist) AddDiagnosis(s string) []string {
	den.Diagnoses = append(den.Diagnoses, s)
	return den.Diagnoses;
}


//TODO: ISSUE with moving from one slice to another !!!!
type patient struct {
	Id string `json:"id"`                        //NHS number in UK
	Name string `json:"name"`
	Address string `json:"address"`
	Open_diagnoses []string `json:"Open"`       //list of open cases
	Closed_diagnoses []string `json:"Close"`    //list of closed cases
}
func (p *patient) AddOpenDiagnosis(s string) []string {
	p.Open_diagnoses = append(p.Open_diagnoses, s)
	return p.Open_diagnoses;
}
func (p *patient) AddClosedDiagnosis(s string) []string {
	p.Closed_diagnoses = append(p.Closed_diagnoses, s)
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
	approved_by []string `json:"approved_by"`   //list of dentists approving the suggested treatment
	disapproved_by []string `json:"disapproved_by"`                     //list of dentists NON-approving the suggested treatment
}
func (d *diagnosis) AddApprovedBy(s string) []string {
	d.approved_by = append(d.approved_by, s)
	return d.approved_by;
}
func (d *diagnosis) AddDisapprovedBy(s string) []string {
	d.disapproved_by = append(d.disapproved_by, s)
	return d.disapproved_by;
}

type MedicalChaincode struct{
}

func (t *MedicalChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Printf("Medical chaincode initialised.\n")
	return nil,nil
}

//functions to add a entities
func (t *MedicalChaincode) addDentist(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) %6 != 0  {
		return nil, errors.New("Incorrect number of args. Needs to be multiples of 6: (Id,Given_name,Surname,Gender,Status,Year_of_qualification")
	}

	DentistOBJ := dentist{}

	for index := 0; index < len(args); index += 6 {
		var id string;
		id = args[index];
		//TODO: look into this duplication, any advantage?
		DentistOBJ.Id = args[index];
		DentistOBJ.Given_name = args[index + 1];
		DentistOBJ.Surname = args[index + 2];
		DentistOBJ.Gender = args[index + 3];
		DentistOBJ.Status = args[index + 4];
		DentistOBJ.Year_of_qualification = args[index + 5];
		var empty []string;
		DentistOBJ.Diagnoses = empty;

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

func (t *MedicalChaincode) addPatient(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) %3 != 0  {
		return nil, errors.New("Incorrect number of args. Needs to be multiples of 3: (Id,name,address")
	}

	PatientOBJ := patient{}

	for index := 0; index < len(args); index += 3 {
		var id string;
		id = args[index];
		//TODO: look into this duplication, any advantage?
		PatientOBJ.Id = args[index];
		PatientOBJ.Name = args[index + 1];
		PatientOBJ.Address = args[index + 2];
		var empty []string;
		PatientOBJ.Open_diagnoses = empty;
		PatientOBJ.Closed_diagnoses = empty;

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

func (t *MedicalChaincode) addDiagnosis(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) %7 != 0  {
		return nil, errors.New("Incorrect number of args. Needs to be multiples of 7: (Id,patient_ID,dentist_ID,date,ICD,treatment,data")
	}

	DiagnosisOBJ := diagnosis{}

	for index := 0; index < len(args); index += 6 {
		var id string;
		id = args[index];
		//TODO: look into this duplication, any advantage?
		DiagnosisOBJ.Id = args[index];
		var patient_ID string;
		patient_ID = args[index + 1];
		DiagnosisOBJ.patient_ID = patient_ID;
		var dentist_ID string;
		dentist_ID = args[index + 2];
		DiagnosisOBJ.dentist_ID = dentist_ID;
		DiagnosisOBJ.Date = args[index + 3];
		DiagnosisOBJ.ICD = args[index + 4];
		DiagnosisOBJ.drug_treatment = args[index + 5];
		DiagnosisOBJ.data = args[index + 6];
		var empty []string;
		DiagnosisOBJ.approved_by = empty;
		DiagnosisOBJ.disapproved_by = empty;

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
		patientOBJ.AddOpenDiagnosis(id);
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
		dentistOBJ.AddDiagnosis(id)
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

func (t *MedicalChaincode) approveDiagnosis(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) %2 != 0  {
		return nil, errors.New("Incorrect number of args. Needs to be multiples of 2: (dentist_ID, diagnosis_ID)")
	}

	diagnosisOBJ := diagnosis{}

	//get the diagnosis by ID
	diagnosisJSON, err := stub.GetState(args[1])
	if diagnosisJSON == nil {
		return nil, errors.New("Error: No account exists for diagnosis.")
	}

	err = json.Unmarshal(diagnosisJSON, &diagnosisOBJ)
	if err != nil {
		return nil, errors.New("Invalid entity data pulled from ledger.")
	}

	diagnosisOBJ.AddApprovedBy(args[0]);

	//TODO: if 3 checks done move to closed or stg in patient


	newdiagnosisJSON, err := json.Marshal(diagnosisOBJ);
	if err != nil || newdiagnosisJSON == nil {
		return nil, errors.New("Converting entity struct to PatientJSON failed")
	}

	//write it back to ledger
	err = stub.PutState(args[1] , newdiagnosisJSON)
	if err != nil {
		fmt.Printf("Error: could not update ledger")
		return nil, err
	}
	return nil,nil
}

func (t *MedicalChaincode) disApproveDiagnosis(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) %2 != 0  {
		return nil, errors.New("Incorrect number of args. Needs to be multiples of 2: (dentist_ID, diagnosis_ID)")
	}

	diagnosisOBJ := diagnosis{}

	//get the diagnosis by ID
	diagnosisJSON, err := stub.GetState(args[1])
	if diagnosisJSON == nil {
		return nil, errors.New("Error: No account exists for diagnosis.")
	}

	err = json.Unmarshal(diagnosisJSON, &diagnosisOBJ)
	if err != nil {
		return nil, errors.New("Invalid entity data pulled from ledger.")
	}

	diagnosisOBJ.AddDisapprovedBy(args[0]);

	//TODO: if 3 checks done move to closed or stg in patient


	newdiagnosisJSON, err := json.Marshal(diagnosisOBJ);
	if err != nil || newdiagnosisJSON == nil {
		return nil, errors.New("Converting entity struct to PatientJSON failed")
	}

	//write it back to ledger
	err = stub.PutState(args[1] , newdiagnosisJSON)
	if err != nil {
		fmt.Printf("Error: could not update ledger")
		return nil, err
	}
	return nil,nil
}

func (t *MedicalChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "addDentist" {
		return t.addDentist(stub, function, args)
	} else if function == "addPatient" {
		return t.addPatient(stub, function, args)
	} else if function == "addDiagnosis" {
		return t.addDiagnosis(stub, function, args)
	}else if function == "approveDiagnosis" {
		return t.approveDiagnosis(stub, function, args)
	}else if function == "disApproveDiagnosis" {
		return t.disApproveDiagnosis(stub, function, args)
	}
	return nil, errors.New("Received unknown function invocation.")
}

func (t *MedicalChaincode) getPatient(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1  {
		return nil, errors.New("Incorrect number of args. Exactly one expected: (patient_ID)")
	}

	patient_ID := args[0]
	dataJson, err := stub.GetState(patient_ID)
	if dataJson == nil || err != nil {
		return nil, errors.New("Cannot get user data from chain.")
	}

	fmt.Printf("Query Response: %s\n", dataJson)
	return dataJson, nil
}


func (t *MedicalChaincode) getDentist(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1  {
		return nil, errors.New("Incorrect number of args. Exactly one expected: (dentist_ID)")
	}

	dentist_ID := args[0]
	dataJson, err := stub.GetState(dentist_ID)
	if dataJson == nil || err != nil {
		return nil, errors.New("Cannot get user data from chain.")
	}

	fmt.Printf("Query Response: %s\n", dataJson)
	return dataJson, nil
}

func (t *MedicalChaincode) getDiagnosis(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1  {
		return nil, errors.New("Incorrect number of args. Exactly one expected: (diagnosis_ID)")
	}

	diagnosis_ID := args[0]
	dataJson, err := stub.GetState(diagnosis_ID)
	if dataJson == nil || err != nil {
		return nil, errors.New("Cannot get user data from chain.")
	}

	fmt.Printf("Query Response: %s\n", dataJson)
	return dataJson, nil
}

func (t *MedicalChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "addDentist" {
		return t.getDentist(stub, function, args)
	} else if function == "addPatient" {
		return t.getPatient(stub, function, args)
	} else if function == "addDiagnosis" {
		return t.getDiagnosis(stub, function, args)
	}
	return nil, errors.New("Received unknown function invocation.")
}


func main() {
	err := shim.Start(new(MedicalChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
