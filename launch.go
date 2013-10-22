package main

import (
  "container/list"
  "log"
  "time"
)

type Job interface {
  Name() string
  Run() error
  Launch() Launch
}

type FuncJob struct {
  closure func() error
  name    string
  launch  Launch
}

func (j *FuncJob) Name() string {
  return j.name
}

func (j *FuncJob) Run() error {
  return j.closure()
}

func (j *FuncJob) Launch() Launch {
  return j.launch
}

type Launch interface {
  NextDuration() *time.Duration
}

type RecurringLaunch time.Duration

func (r *RecurringLaunch) NextDuration() *time.Duration {
  return (*time.Duration)(r)
}

type Entry struct {
  Job Job
}

type Launcher struct {
  Jobs     []Job
  JobQueue list.List
}

func (l *Launcher) Schedule(j Job) error {
  l.JobQueue.PushBack(j)
  return nil
}

func (l *Launcher) Process() error {
  for {
    job_elt := l.JobQueue.Front()
    if job_elt == nil {
      break
    }

    l.JobQueue.Remove(job_elt)
    job := job_elt.Value.(Job)
    dur := job.Launch().NextDuration()
    if dur == nil {
      break
    }
    time.Sleep(*dur)
    go job.Run()
    l.JobQueue.PushBack(job)
  }
  return nil
}

func main() {
  duration, err := time.ParseDuration("5s")
  if err != nil {
    log.Fatal(err)
  }

  job := &FuncJob{
    name: "Hello job",
    closure: func() error {
      println("Hello, world.")
      return nil
    },
    launch: (*RecurringLaunch)(&duration),
  }
  launcher := new(Launcher)
  launcher.Schedule(job)
  launcher.Process()
}
