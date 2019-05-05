package client

import (
	"encoding/json"
	"fossil/src/engine/util"
	"net/http"
	"bytes"
	"errors"
)

func PreQuiesceCmd(auth Auth,config util.Config) (util.Result,error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-app:8001/preQuiesceCmd", b)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func QuiesceCmd(auth Auth,config util.Config) (util.Result,error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-app:8001/quiesceCmd", b)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func PostQuiesceCmd(auth Auth,config util.Config) (util.Result,error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-app:8001/postQuiesceCmd", b)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func PreUnquiesceCmd(auth Auth,config util.Config) (util.Result,error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-app:8001/preUnquiesceCmd", b)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func UnquiesceCmd(auth Auth,config util.Config) (util.Result,error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-app:8001/unquiesceCmd", b)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func BackupCreateCmd(auth Auth,config util.Config) (util.Result,error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/backupCreateCmd", b)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func BackupDeleteCmd(auth Auth,config util.Config) (util.Result,error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/backupDeleteCmd", b)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func PostUnquiesceCmd(auth Auth,config util.Config) (util.Result,error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-app:8001/postUnquiesceCmd", b)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}


func SendTrapSuccessCmd(auth Auth,config util.Config) (util.Result,error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-workflow:8000/sendTrapSuccessCmd", b)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func SendTrapErrorCmd(auth Auth,config util.Config) (util.Result,error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-workflow:8000/sendTrapErrorCmd", b)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func ArchiveCreateCmd(auth Auth,config util.Config) (util.Result,error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/archiveCreateCmd", b)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func ArchiveDeleteCmd(auth Auth,config util.Config) (util.Result,error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/archiveDeleteCmd", b)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}