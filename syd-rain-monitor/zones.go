package main



type Zone struct{
	Name string
	Offset int
	Threshold int
	Score int
}

func (z *Zone)Triggered()bool{
	return z.Score > z.Threshold
}

func (z *Zone)Percentage()int{
	return int(float32(z.Score) / float32(z.Threshold) * 100)
}

var zones = []*Zone{
	&Zone{
		Name: "Center",
		Offset: 6,
		Threshold: 55,
	},
	&Zone{
		Name: "Imminent",
		Offset: 16,
		Threshold: 200,
	},
	&Zone{
		Name: "Close",
		Offset: 27,
		Threshold: 200,
	},
	&Zone{
		Name: "Warning",
		Offset: 37,
		Threshold: 200,
	},
	&Zone{
		Name: "Outer",
		Offset: 52,
		Threshold: 200,
	},
}
