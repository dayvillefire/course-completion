package main

import (
	"flag"
	"log"
	"strings"
)

var (
	action        = flag.String("action", "", "Action (test|run|dryrun)")
	configFile    = flag.String("config", "config.yaml", "YAML configuration file")
	dataFile      = flag.String("data", "in.xslx", "Input excel spreadsheet")
	dataSheetName = flag.String("data-sheet-name", "Sheet1", "Excel sheet name")
)

func main() {
	flag.Parse()

	// Load the configuration for the job
	c, err := loadConfig(*configFile)
	if err != nil {
		panic(err)
	}

	// Determine action
	if *action != "test" && *action != "run" && *action != "dryrun" {
		panic("Invalid action")
	}

	// If we're testing, generate a single test email with a cert
	// then exit
	if *action == "test" {
		// Get the list of recipients from a spreadsheet

		// Iterate:

		destName := "Jeffrey Buchbinder"
		destEmail := "j.buchbinder@dayvillefire.gov"
		certName := "Certificate - " + destName + ".pdf"

		// Create certificate
		log.Printf("%#v", c)
		g := generator{
			Template:           c.Template.PdfFile,
			GlobalReplacements: c.Template.GlobalReplacements,
			Replacements:       c.Template.Replacements,
		}

		defaults := map[string]string{}
		for x, y := range c.Template.Replacements {
			defaults[x] = y.Default
		}

		err := g.generate(certName, defaults)
		if err != nil {
			panic(err)
		}

		// Email certificate
		m := mailer{
			MailServer: c.Mail.ServerName,
			Port:       c.Mail.ServerPort,
			Username:   c.Mail.Username,
			Password:   c.Mail.Password,
		}
		err = m.sendMail(
			destName,
			destEmail,
			c.Mail.FromName,
			c.Mail.FromEmail,
			c.Mail.Subject,
			c.Mail.Template,
			certName,
		)
		if err != nil {
			panic(err)
		}

		// Print some sort of acknowledgement

		return
	}

	// If we're doing a full run:
	if *action == "run" || *action == "dryrun" {

		// Get the list of recipients from a spreadsheet
		data, err := ReadExcel(*dataFile, *dataSheetName)

		// Iterate:
		for iter, item := range data {
			log.Printf("INFO: Iteration #%d", iter+1)

			destName, ok := item[c.Data.NameField]
			if !ok {
				log.Printf("WARN: Skipped item, no data for %s column", c.Data.NameField)
				continue
			}
			destEmail, ok := item[c.Data.EmailField]
			if !ok {
				log.Printf("WARN: Skipped item, no data for %s column", c.Data.EmailField)
				continue
			}
			certName := "Certificate - " + strings.TrimSpace(destName) + ".pdf"

			// Create certificate
			log.Printf("INFO: %#v", c)
			g := generator{
				Template:     c.Template.PdfFile,
				Replacements: c.Template.Replacements,
			}

			values := map[string]string{}

			// Populate defaults
			for x, y := range c.Template.Replacements {
				values[x] = y.Default
			}

			// Override values with global replacements
			for k, v := range c.Template.GlobalReplacements {
				values[k] = v
			}

			// Override values with data
			for k, v := range item {
				values[k] = v
			}

			err = g.generate(certName, values)
			if err != nil {
				log.Printf("ERR: %s", err.Error())
				continue
			}

			if *action != "dryrun" {
				// Email certificate
				m := mailer{
					MailServer: c.Mail.ServerName,
					Port:       c.Mail.ServerPort,
					Username:   c.Mail.Username,
					Password:   c.Mail.Password,
				}
				err = m.sendMail(
					destName,
					destEmail,
					c.Mail.FromName,
					c.Mail.FromEmail,
					c.Mail.Subject,
					c.Mail.Template,
					certName,
				)
				if err != nil {
					log.Printf("ERR: %s", err.Error())
				}
			}

			log.Printf("INFO: Sent certificate to %s <%s>", destName, destEmail)
		}

		// Print some sort of acknowledgement
	}
}
