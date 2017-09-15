/*
Copyright 2017 Rohith Jayawardene All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:    "cfssl-sidekick",
		Author:  "Rohith Jayawardene",
		Email:   "gambol99@gmail.com",
		Usage:   "is a small utilty service used to acquire certificates from cfssl",
		Version: fmt.Sprintf("%s (git+sha: %s)", Version, GitSHA),

		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "cfssl-token",
				Usage:  "a authentication token for cfssl `TOKEN`",
				EnvVar: "CFSSL_TOKEN",
			},
			cli.StringFlag{
				Name:   "cfssl-url",
				Usage:  "a cfssl endpoint url for the service `URL`",
				EnvVar: "CFSSL_URL",
				Value:  "https://ca.kube-tls.svc.cluster.local",
			},
			cli.StringFlag{
				Name:   "cfssl-profile",
				Usage:  "a cfssl profile to use when requesting a certificated `NAME`",
				EnvVar: "CFSSL_PROFILE",
				Value:  "default",
			},
			cli.StringFlag{
				Name:   "tls-ca",
				Usage:  "the path to a file containing tls certificate for CA `PATH`",
				EnvVar: "TLS_CA",
			},
			cli.StringSliceFlag{
				Name:   "domain",
				Usage:  "a list of domains you are requesting for `DOMAIN`",
				EnvVar: "DOMAIN",
			},
			cli.IntFlag{
				Name:   "size",
				Usage:  "the size of the certificate `SIZE`",
				Value:  2048,
				EnvVar: "SIZE",
			},
			cli.StringFlag{
				Name:   "output-dir",
				Usage:  "the path to the directory where the certificates should be saved `PATH`",
				Value:  "/certs",
				EnvVar: "OUTPUT_DIR",
			},
			cli.StringFlag{
				Name:   "command",
				Usage:  "an command line executable to run when a new certificate is acquired `COMMAND`",
				EnvVar: "COMMAND",
			},
			cli.StringFlag{
				Name:   "tls-organization",
				Usage:  "the organization name for the certificate request `NAME`",
				Value:  "ACP Homeoffice",
				EnvVar: "TLS_ORGANIZATION",
			},
			cli.StringFlag{
				Name:   "tls-country",
				Usage:  "the country name placed into the certificate request `NAME`",
				Value:  "GB",
				EnvVar: "TLS_COUNTRY",
			},
			cli.StringFlag{
				Name:   "tls-locality",
				Usage:  "the locality name placed into the certificate `NAME`",
				Value:  "London",
				EnvVar: "TLS_LOCALITY",
			},
			cli.StringFlag{
				Name:   "tls-province",
				Usage:  "the province name placed in the certificate `NAME`",
				Value:  "London",
				EnvVar: "TLS_PROVINCE",
			},
			cli.BoolFlag{
				Name:   "onetime",
				Usage:  "indicated you only want the service to run once and exit `BOOL`",
				EnvVar: "ONETIME",
			},
			cli.DurationFlag{
				Name:   "timeout",
				Usage:  " a timeout for operation, if we've not recieved a certificate in this time, exit `DURATION`",
				Value:  1 * time.Minute,
				EnvVar: "TIMEOUT",
			},
			cli.BoolFlag{
				Name:   "verbose",
				Usage:  "whether to enable verbose logging `BOOL`",
				EnvVar: "VERBOSE",
			},
		},

		Action: func(c *cli.Context) error {
			cfg := Config{
				CertsDir:        c.String("output-dir"),
				Country:         c.String("tls-country"),
				Domains:         c.StringSlice("domain"),
				EndpointProfile: c.String("cfssl-token"),
				EndpointToken:   c.String("cfssl-token"),
				EndpointURL:     c.String("cfssl-url"),
				ExecCommand:     c.String("exec"),
				Locality:        c.String("tls-locality"),
				Onetime:         c.Bool("onetime"),
				Organization:    c.String("tls-organization"),
				Province:        c.String("tls-province"),
				Size:            c.Int("size"),
				Timeout:         c.Duration("timeout"),
				TLSCAPath:       c.String("tls-ca"),
				Verbose:         c.Bool("verbose"),
			}

			ctl, err := newController(cfg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[error] failed to initialize controller, error: %s\n", err)
				os.Exit(1)
			}

			if err := ctl.run(); err != nil {
				fmt.Fprintf(os.Stderr, "[error] failed to start controller, error: %s\n", err)
				os.Exit(1)
			}

			return nil
		},
	}

	app.Run(os.Args)
}
