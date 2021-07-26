package data

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/buntdb"
)

// UpdateDefaultUSBDev updates the default usb device to use
func (store Manager) UpdateDefaultUSBDev(updatedDefaultDev string) (string, error) {

	//	Our return item
	retval := ""

	//	Serialize to JSON format
	encoded, err := json.Marshal(updatedDefaultDev)
	if err != nil {
		return retval, fmt.Errorf("problem serializing the data: %s", err)
	}

	//	Save it to the database:
	err = store.systemdb.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(GetKey("Config", "DefaultUSBDevice"), string(encoded), &buntdb.SetOptions{})
		return err
	})

	//	If there was an error saving the data, report it:
	if err != nil {
		return retval, fmt.Errorf("problem saving the default USB device: %s", err)
	}

	//	Set our retval:
	retval = updatedDefaultDev

	//	Return our data:
	return retval, nil
}

// GetDefaultUSBDev gets the default usb device to use.  If not set, it attempt to set it first.  Returns an error if it can't be fetched or set
func (store Manager) GetDefaultUSBDev() (string, error) {
	//	Our return item
	retval := ""

	//	Find the item:
	err := store.systemdb.View(func(tx *buntdb.Tx) error {

		val, err := tx.Get(GetKey("Config", "DefaultUSBDevice"))
		if err != nil {
			return err
		}

		if len(val) > 0 {
			//	Unmarshal data into our item
			if err := json.Unmarshal([]byte(val), &retval); err != nil {
				return err
			}
		}

		//	If we get to this point and there is no error...
		return nil
	})

	//	If there was an error, report it:
	if err != nil {
		return retval, fmt.Errorf("problem getting the default device: %s", err)
	}

	//	Return our data:
	return retval, nil
}
