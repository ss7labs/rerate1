package main

import (
	//"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Job struct {
	orgId int
	date  string
	dual  bool
}

func runWork(id int, date string, dual bool, reestr *Reestr) {
	var idStr = strconv.Itoa(id)
	var org Org
	org.Id = id
	org.Date = date
	org.dual = dual
	org.reestr = reestr
	query := "SELECT org_name,org_acct,firma FROM orgs WHERE org_id=" + idStr
	err := asudb.QueryRow(query).Scan(&org.Name, &org.Acct, &org.Firma)

	if err != nil {
		if err == sql.ErrNoRows {
			return
		} else {
			panic(err.Error())
		}
	}

	if org.dual {
		org.printDualOrg()
	} else {
		org.printOneOrg()
	}

	log.WithFields(logrus.Fields{
		"id": id,
	}).Debug("runWork")
}

func getJobsDual(date string) []Job {
	query := "SELECT org_id FROM orgs WHERE org_id IN (3095,2481,2245,3011) ORDER BY org_id ASC"
	return getJobsHelper(query, date, true)
}

func getJobs(date string) []Job {
	/* Akhal Orgs*/
	//query := "SELECT org_id FROM orgs WHERE org_acct LIKE '709%' ORDER BY org_id ASC"
	/* MVD Orgs*/
	query := "SELECT org_id FROM orgs WHERE org_acct LIKE '280317%' ORDER BY org_id ASC"
	/* All Orgs*/
	// query := "SELECT org_id FROM orgs ORDER BY org_id ASC"
	/*
	    query := "SELECT org_id FROM orgs WHERE org_id NOT IN (3095,2481,2245,3011) ORDER BY org_id ASC"
	    query := "SELECT org_id FROM orgs WHERE org_acct LIKE '204256%' ORDER BY org_id ASC"
	    query := `SELECT org_id FROM orgs WHERE org_id IN (
	   -165,
	   -166,
	   709,
	   1217,
	   838,
	   -171,
	   8122,
	   -172,
	   685,
	   164,
	   2720,
	   2167,
	   3571,
	   -170,
	   1921,
	   -2153,
	   3166,
	   3351,
	   -168,
	   1885,
	   2472,
	   -167,
	   4540,
	   -169,
	   2288,
	   -173,
	   -128,
	   591,
	   -174) ORDER BY org_id ASC`
	*/
	return getJobsHelper(query, date, false)
}

func getJobsHelper(query, date string, dual bool) []Job {

	var jobs []Job
	rows, err := asudb.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var orgId int
		if err := rows.Scan(&orgId); err != nil {
			panic(err.Error())
		}
		jobs = append(jobs, Job{orgId, date, dual})
		//fmt.Println("getJobs",orgId)
	}

	return jobs

}
