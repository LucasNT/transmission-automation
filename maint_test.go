
package main

import (
    "testing"
)


func TestReadConfigsFromYaml( t *testing.T) {
    expected := Config{"senha_super_secreto", "http://url_do_transmission/transmission/rpc/", "usuário"}
    actual,err := readConfigsFromYaml("./testdata/config_test.yaml");
    if err != nil {
        t.FailNow();
    }
    if expected != actual {
        t.Errorf("Expected %s do not match actual %s", expected, actual)
    }
}
