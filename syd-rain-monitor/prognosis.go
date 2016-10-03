package main


type IProgEngine interface{
  Name()string
  Feed([]*Zone)
  Good()bool
  Certainity()int
  CertainityAvailable()bool
}



type ZoningPrognosis struct{
  latestZones []*Zone

  lastNumTriggered int
  currentNumTriggered int

  lastTimestamp string
  currentTimestamp string
}

func (p *ZoningPrognosis)Name()string{
  return "Zoning"
}

func (p *ZoningPrognosis)Feed(zones []*Zone){
  p.latestZones = zones

  p.lastTimestamp = p.currentTimestamp
  p.currentTimestamp = updated

  p.lastNumTriggered = p.currentNumTriggered
  p.currentNumTriggered = 0
  for _, zone := range p.latestZones {
    if zone.Triggered() {
      p.currentNumTriggered += 1
    }
  }
}

func (p *ZoningPrognosis)Good()bool{
  if p.currentNumTriggered >= 1 && p.lastNumTriggered >= 1 && p.lastTimestamp != p.currentTimestamp {
    return false
  }

  return true
}

func (p *ZoningPrognosis)Certainity()int{
  certainity := 100

  for _, zone := range p.latestZones {
    if zone.Triggered() {
      certainity -= (100/(len(p.latestZones)+1))
    }
  }

  return certainity-5
}

func (p *ZoningPrognosis)CertainityAvailable()bool{
  return true
}
