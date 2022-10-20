package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
	"github.com/pja237/slurmcommander/internal/keybindings"
	"github.com/pja237/slurmcommander/internal/styles"
)

// genTabs() generates top tabs
func (m Model) genTabs() string {

	var doc strings.Builder

	tlist := make([]string, len(tabs))
	for i, v := range tabs {
		if i == int(m.ActiveTab) {
			tlist = append(tlist, styles.TabActiveTab.Render(v))
		} else {
			tlist = append(tlist, styles.Tab.Render(v))
		}
	}
	row := lipgloss.JoinHorizontal(lipgloss.Top, tlist...)

	//gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
	gap := styles.TabGap.Render(strings.Repeat(" ", max(0, m.winW-lipgloss.Width(row)-2)))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
	doc.WriteString(row + "\n")

	return doc.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m Model) tabJobs() string {

	// TODO: See what's more visually clear way to present info
	// e.g. Show selected job info in:
	// a) separate toggle-able window on the side OR
	// b) header/footer (above/below) of the table, like `top` does
	return m.SqueueTable.View() + "\n"
}

func (m Model) tabJobHist() string {

	// TODO: do some statistics on job history
	// e.g. avg waiting times, jobs successfull/failed count, etc...

	return m.JobHistTab.SacctTable.View() + "\n"
}

func (m Model) tabJobDetails() (scr string) {

	// race between View() call and command.SingleJobGetSacct(m.JobDetailsTab.SelJobID) call
	switch {
	case m.JobDetailsTab.SelJobID == "":
		return "Select a job from the Job History tab.\n"
	case len(m.SacctSingleJobHist.Jobs) == 0:
		return fmt.Sprintf("Waiting for job %s info...\n", m.JobDetailsTab.SelJobID)

	}

	width := m.Globals.winW - 10
	job := m.SacctSingleJobHist.Jobs[0]

	m.Log.Printf("Job Details req %#v ,got: %#v\n", m.JobDetailsTab.SelJobID, job.JobId)
	//scr = fmt.Sprintf("Job count: %d\n\n", len(m.SacctJob.Jobs))

	// TODO: consider moving this to a table...

	head := ""
	waitT := time.Unix(int64(*job.Time.Start), 0).Sub(time.Unix(int64(*job.Time.Submission), 0))
	runT := time.Unix(int64(*job.Time.End), 0).Sub(time.Unix(int64(*job.Time.Start), 0))
	fmtStr := "%-20s : %-40s\n"
	head += fmt.Sprintf(fmtStr, "Job ID", strconv.Itoa(*job.JobId))
	head += fmt.Sprintf(fmtStr, "User", *job.User)
	head += fmt.Sprintf(fmtStr, "Job Name", *job.Name)
	head += fmt.Sprintf(fmtStr, "Job Account", *job.Account)
	head += fmt.Sprintf(fmtStr, "Job Submission", time.Unix(int64(*job.Time.Submission), 0).String())
	head += fmt.Sprintf(fmtStr, "Job Start", time.Unix(int64(*job.Time.Start), 0).String())
	head += fmt.Sprintf(fmtStr, "Job End", time.Unix(int64(*job.Time.End), 0).String())
	head += fmt.Sprintf(fmtStr, "Job Wait time", waitT.String())
	head += fmt.Sprintf(fmtStr, "Job Run time", runT.String())
	head += fmt.Sprintf(fmtStr, "Partition", *job.Partition)
	head += fmt.Sprintf(fmtStr, "Priority", strconv.Itoa(*job.Priority))
	head += fmt.Sprintf(fmtStr, "QoS", *job.Qos)

	scr += styles.JobStepBoxStyle.Width(width).Render(head)
	scr += fmt.Sprintf("\n Steps count: %d", len(*job.Steps))

	steps := ""
	for i, v := range *job.Steps {

		m.Log.Printf("Job Details, step: %d name: %s\n", i, *v.Step.Name)
		step := fmt.Sprintf(fmtStr, "Name", *v.Step.Name)
		step += fmt.Sprintf(fmtStr, "State", *v.State)
		step += fmt.Sprintf(fmtStr, "ExitStatus", *v.ExitCode.Status)
		if *v.ExitCode.Status == "SIGNALED" {
			step += fmt.Sprintf(fmtStr, "Signal ID", strconv.Itoa(*v.ExitCode.Signal.SignalId))
			step += fmt.Sprintf(fmtStr, "SignalName", *v.ExitCode.Signal.Name)
		}
		if v.KillRequestUser != nil {
			step += fmt.Sprintf(fmtStr, "KillReqUser", *v.KillRequestUser)
		}
		step += fmt.Sprintf(fmtStr, "Tasks", strconv.Itoa(*v.Tasks.Count))
		//steps += lipgloss.JoinVertical(lipgloss.Bottom, steps, styles.JobStepBoxStyle.Width(m.Globals.winW-10).Render(step))
		steps += "\n" + styles.JobStepBoxStyle.Width(width).Render(step)
		//m.Log.Printf("Step %d, VALUE: %q", i, steps)
		//m.Log.Printf("================================================================================\n")
	}
	scr += steps

	//// get Type of struct
	//t := reflect.TypeOf(m.JobDetailsTab.Jobs[0])
	//// get Value of struct
	//v := reflect.ValueOf(m.JobDetailsTab.Jobs[0])
	//// get struct fields from Type
	//f := reflect.VisibleFields(t)
	//scr += fmt.Sprintf("num. Fields = %d\n", len(f))
	//for _, val := range f {
	//	dv := v.FieldByName(val.Name).Elem()
	//	switch dv.Kind() {
	//	case reflect.String:
	//		scr += fmt.Sprintf("*string FIELD: %-20s VALUE: %-20s\n", val.Name, dv.String())
	//	case reflect.Int:
	//		scr += fmt.Sprintf("*int FIELD: %-20s VALUE: %-20d\n", val.Name, dv.Int())
	//	default:
	//		scr += fmt.Sprintf("UnknownKind %v FIELD: %-20s VALUE: %-20v\n", v.FieldByName(val.Name).Kind(), val.Name, v.FieldByName(val.Name).String())

	//	}
	//}
	//scr += fmt.Sprintf("Job:\n\n%#v\n\nSelected job: %#v\n\n", m.JobDetailsTab.SacctJob, m.JobDetailsTab.SelJobID)
	//m.LogF.WriteString(fmt.Sprintf("Job:\n\n%#v\n\nSelected job: %#v\n\n", m.JobDetailsTab.SacctJob, m.JobDetailsTab.SelJobID))

	return scr
	//return m.JobDetailsTab.SelJobID
}

func (m Model) tabJobFromTemplate() string {

	if m.EditTemplate {
		return m.TemplateEditor.View()
	} else {
		// TODO: if len(table)==0 return "no templates found"
		if len(m.JobFromTemplateTab.TemplatesList) == 0 {
			return styles.NotFound.Render("\nNo templates found!\n")
		} else {
			return m.TemplatesTable.View()
		}
	}
}

func (m Model) tabCluster() string {

	var (
		scr     string = ""
		cpuPerc float64
		memPerc float64
	)

	// node info
	// TODO: rework, doesn't work when table filtering is on
	sel := m.SinfoTable.Cursor()
	m.Log.Printf("ClusterTab Selected: %d\n", sel)
	m.Log.Printf("ClusterTab len results: %d\n", len(m.JobClusterTab.SinfoFiltered.Nodes))
	if len(m.JobClusterTab.SinfoFiltered.Nodes) > 0 && sel != -1 {
		m.CpuBar = progress.New(progress.WithGradient("#277BC0", "#FFCB42"))
		cpuPerc = float64(*m.JobClusterTab.SinfoFiltered.Nodes[sel].AllocCpus) / float64(*m.JobClusterTab.SinfoFiltered.Nodes[sel].Cpus)
		//m.CpuBar.SetPercent(cpuPerc)
		m.MemBar = progress.New(progress.WithGradient("#277BC0", "#FFCB42"))
		memPerc = float64(*m.JobClusterTab.SinfoFiltered.Nodes[sel].AllocMemory) / float64(*m.JobClusterTab.SinfoFiltered.Nodes[sel].RealMemory)
	} else {
		cpuPerc = 0
		memPerc = 0
	}

	scr += "Cpu and memory utilization:\n"
	scr += fmt.Sprintf("cpuPerc: %.2f ", cpuPerc)
	scr += m.CpuBar.ViewAs(cpuPerc)
	scr += "\n"
	scr += fmt.Sprintf("memPerc: %.2f ", memPerc)
	scr += m.MemBar.ViewAs(memPerc)
	scr += "\n\n"

	// table
	scr += m.SinfoTable.View() + "\n"

	return scr
	//return m.SinfoTable.View()

	//return "Cluster tab active"
	// TODO: HEADER/FOOTER that shows details for selected node
	// e.g. progress bars with cpu/mem usage (percentages)
	// TODO: reselect table columns (move mem/cpu to header/footer above and pick others, e.g. partition? features? think... )
	// TODO: rename to "Cluster Nodes", add "Cluster QoS/Partition" tab? OR find an elegant way to group those in one tab?
}

func (m Model) tabAbout() string {

	s := `

petar.jager@imba.oeaw.ac.at
CLIP-HPC Team @ VBC
	`

	return "About tab active" + s
}

func (m Model) getJobInfo() string {
	var scr strings.Builder

	n := m.JobTab.SqueueTable.Cursor()
	m.Log.Printf("getJobInfo: cursor at %d table rows: %d\n", n, len(m.JobTab.SqueueFiltered.Jobs))
	if len(m.JobTab.SqueueFiltered.Jobs) == 0 || n == -1 {
		return "Select a job"
	}

	fmtStr := "%-15s : %-30s\n"
	fmtStrLast := "%-15s : %-30s"
	//ibFmt := "Job Name: %s\nJob Command: %s\nOutput: %s\nError: %s\n"
	infoBoxLeft := fmt.Sprintf(fmtStr, "Partition", *m.JobTab.SqueueFiltered.Jobs[n].Partition)
	infoBoxLeft += fmt.Sprintf(fmtStr, "QoS", *m.JobTab.SqueueFiltered.Jobs[n].Qos)
	infoBoxLeft += fmt.Sprintf(fmtStr, "TRES", *m.JobTab.SqueueFiltered.Jobs[n].TresReqStr)
	if m.JobTab.SqueueFiltered.Jobs[n].JobResources.Nodes != nil {
		infoBoxLeft += fmt.Sprintf(fmtStr, "AllocNodes", *m.JobTab.SqueueFiltered.Jobs[n].JobResources.Nodes)
	} else {
		infoBoxLeft += fmt.Sprintf(fmtStr, "AllocNodes", "none")

	}
	infoBoxLeft += fmt.Sprintf(fmtStrLast, "wckey", *m.JobTab.SqueueFiltered.Jobs[n].Wckey)

	infoBoxRight := fmt.Sprintf(fmtStr, "Array Job ID", strconv.Itoa(*m.JobTab.SqueueFiltered.Jobs[n].ArrayJobId))
	if m.JobTab.SqueueFiltered.Jobs[n].ArrayTaskId != nil {
		infoBoxRight += fmt.Sprintf(fmtStr, "Array Task ID", strconv.Itoa(*m.JobTab.SqueueFiltered.Jobs[n].ArrayTaskId))
	} else {
		infoBoxRight += fmt.Sprintf(fmtStr, "Array Task ID", "NoTaskID")
	}
	infoBoxRight += fmt.Sprintf(fmtStr, "Gres Details", strings.Join(*m.JobTab.SqueueFiltered.Jobs[n].GresDetail, ","))
	infoBoxRight += fmt.Sprintf(fmtStr, "Batch Host", *m.JobTab.SqueueFiltered.Jobs[n].BatchHost)
	infoBoxRight += fmt.Sprintf(fmtStrLast, "Features", *m.JobTab.SqueueFiltered.Jobs[n].Features)

	infoBoxMiddle := fmt.Sprintf(fmtStr, "Submit", time.Unix(*m.JobTab.SqueueFiltered.Jobs[n].SubmitTime, 0))
	if *m.JobTab.SqueueFiltered.Jobs[n].StartTime != 0 {
		infoBoxMiddle += fmt.Sprintf(fmtStrLast, "Start", time.Unix(*m.JobTab.SqueueFiltered.Jobs[n].StartTime, 0))
	} else {
		infoBoxMiddle += fmt.Sprintf(fmtStrLast, "Start", "unknown")
	}

	infoBoxWide := fmt.Sprintf(fmtStr, "Job Name", *m.JobTab.SqueueFiltered.Jobs[n].Name)
	infoBoxWide += fmt.Sprintf(fmtStr, "Command", *m.JobTab.SqueueFiltered.Jobs[n].Command)
	infoBoxWide += fmt.Sprintf(fmtStr, "StdOut", *m.JobTab.SqueueFiltered.Jobs[n].StandardOutput)
	infoBoxWide += fmt.Sprintf(fmtStr, "StdErr", *m.JobTab.SqueueFiltered.Jobs[n].StandardError)
	infoBoxWide += fmt.Sprintf(fmtStrLast, "Working Dir", *m.JobTab.SqueueFiltered.Jobs[n].CurrentWorkingDirectory)

	top := lipgloss.JoinHorizontal(lipgloss.Top, styles.JobInfoInBox.Render(infoBoxLeft), styles.JobInfoInBox.Render(infoBoxMiddle), styles.JobInfoInBox.Render(infoBoxRight))
	scr.WriteString(lipgloss.JoinVertical(lipgloss.Left, top, styles.JobInfoInBox.Render(infoBoxWide)))

	//return infoBox
	return scr.String()
}

func genTabHelp(t int) string {
	var th string
	switch t {
	case tabJobs:
		th = "List of jobs in the queue"
	case tabJobHist:
		th = "Last 7 days job history"
	case tabJobDetails:
		th = "Job details, select a job from Job History tab"
	case tabJobFromTemplate:
		th = "Edit and submit one of the job templates"
	case tabCluster:
		th = "List and status of cluster nodes"
	default:
		th = "SlurmCommander"
	}
	return th + "\n\n"
}

func (m Model) View() string {

	var scr strings.Builder

	// HEADER / TABS
	scr.WriteString(m.genTabs())
	scr.WriteString(genTabHelp(int(m.ActiveTab)))

	// PICK and RENDER ACTIVE TAB

	switch m.ActiveTab {
	case tabJobs:
		//scr.WriteString("Filter: " + m.JobTab.Filter.Value() + "\n\n")
		scr.WriteString(fmt.Sprintf("Filter: %10.10s\tItems: %d\n", m.JobTab.Filter.Value(), len(m.JobTab.SqueueFiltered.Jobs)))
		scr.WriteString("Count: ")
		for k, v := range m.JobTab.Stats.StateCnt {
			scr.WriteString(fmt.Sprintf("%s: %d ", k, v))
		}
		scr.WriteString("\n\n")

		switch {
		case m.FilterSwitch == FilterSwitch(m.ActiveTab):
			scr.WriteString(m.tabJobs())
			scr.WriteString(fmt.Sprintf("Filter value (search accross all fields!):\n%s\n%s", m.JobTab.Filter.View(), "(Enter to finish, Esc to clear filter and abort)") + "\n")
		case m.JobTab.MenuOn:
			// TODO: Render menu here
			scr.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, m.tabJobs(), styles.MenuBoxStyle.Render(m.JobTab.Menu.View())))
			//scr.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, m.tabJobs(), m.JobTab.Menu.View()))
			m.Log.Printf("\nITEMS LIST: %#v\n", m.JobTab.Menu.Items())
		case m.JobTab.InfoOn:
			//scr.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, m.tabJobs(), styles.JobInfoBox.Render(m.getJobInfo())))
			scr.WriteString(m.tabJobs() + "\n")
			scr.WriteString(styles.JobInfoBox.Render(m.getJobInfo()))
		default:
			scr.WriteString(m.tabJobs())
		}
	case tabJobHist:
		//scr.WriteString("Filter: " + m.JobHistTab.Filter.Value() + "\n\n")
		scr.WriteString(fmt.Sprintf("Filter: %10.10s\tItems: %d\n\n", m.JobHistTab.Filter.Value(), len(m.JobHistTab.SacctHistFiltered.Jobs)))
		switch {
		case m.FilterSwitch == FilterSwitch(m.ActiveTab):
			scr.WriteString(m.tabJobHist())
			scr.WriteString(fmt.Sprintf("Filter value (search accross all fields!):\n%s\n%s", m.JobHistTab.Filter.View(), "(Enter to finish, Esc to clear filter and abort)") + "\n")
		default:
			scr.WriteString(m.tabJobHist())
		}
	case tabJobDetails:
		scr.WriteString(m.tabJobDetails())
	case tabJobFromTemplate:
		scr.WriteString(m.tabJobFromTemplate())
	case tabCluster:
		//scr.WriteString("Filter: " + m.JobClusterTab.Filter.Value() + "\n\n")
		scr.WriteString(fmt.Sprintf("Filter: %10.10s\tItems: %d\n\n", m.JobClusterTab.Filter.Value(), len(m.JobClusterTab.SinfoFiltered.Nodes)))
		switch {
		case m.FilterSwitch == FilterSwitch(m.ActiveTab):
			scr.WriteString(m.tabCluster())
			scr.WriteString(fmt.Sprintf("Filter value (search accross all fields!):\n%s\n%s", m.JobClusterTab.Filter.View(), "(Enter to finish, Esc to clear filter and abort)") + "\n")
		default:
			scr.WriteString(m.tabCluster())
		}
	case tabAbout:
		scr.WriteString(m.tabAbout())
	}

	// FOOTER
	scr.WriteString("\n")
	// Debug information:
	//if m.Globals.Debug {
	//	scr.WriteString("DEBUG:\n")
	//	scr.WriteString(fmt.Sprintf("Last key pressed: %q\n", m.lastKey))
	//	scr.WriteString(fmt.Sprintf("Window Width: %d\tHeight:%d\n", m.winW, m.winH))
	//	scr.WriteString(fmt.Sprintf("Active tab: %d\t Active Filter value: TBD\t InfoOn: %v\n", m.ActiveTab, m.InfoOn))
	//	scr.WriteString(fmt.Sprintf("Debug Msg: %q\n", m.DebugMsg))
	//}

	// TODO: Help doesn't split into multiple lines (e.g. when window too narrow)
	scr.WriteString(m.Help.View(keybindings.DefaultKeyMap))

	return scr.String()
}
