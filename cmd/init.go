package cmd

import (
	"strings"

	"github.com/jbirtley88/gremel/util"
	log "github.com/sirupsen/logrus"
)

func initialise() {
	if cfgFile != "" {
		err := util.LoadConfig(cfgFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Override any config from the command-line with '--set foo.bar=123 ...'
	for _, kv := range configOverrides {
		fields := strings.Split(kv, "=")
		if len(fields) < 2 {
			log.Fatalf("'--set' needs a 'name=VALUE' argument")
		}
		value := fields[1]
		if len(fields) > 2 {
			value = strings.Join(fields[1:len(fields)-1], "=")
		}
		log.Infof("Setting config: %s = %s", fields[0], value)
	}
}
