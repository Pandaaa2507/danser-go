package difficulty

import "math"

const (
	HitFadeIn     = 400
	HitFadeOut    = 240
	HittableRange = 400
	ResultFadeIn  = 120
	ResultFadeOut = 600
	PostEmpt      = 500
)

type Difficulty struct {
	hpDrain, cs, od, ar float64
	Preempt, FadeIn     float64
	CircleRadius        float64
	Mods                Modifier
	Hit50               int64
	Hit100              int64
	Hit300              int64
	HPMod               float64
	SpinnerRatio        float64
}

func NewDifficulty(hpDrain, cs, od, ar float64) *Difficulty {
	diff := new(Difficulty)
	diff.hpDrain = hpDrain
	diff.cs = cs
	diff.od = od
	diff.ar = ar
	diff.calculate()
	return diff
}

func (diff *Difficulty) calculate() {
	hpDrain, cs, od, ar := diff.hpDrain, diff.cs, diff.od, diff.ar

	if diff.Mods&HardRock > 0 {
		ar = math.Min(ar*1.4, 10)
		cs = math.Min(cs*1.3, 10)
		od = math.Min(od*1.4, 10)
		hpDrain = math.Min(hpDrain*1.4, 10)
	}

	if diff.Mods&Easy > 0 {
		ar /= 2
		cs /= 2
		od /= 2
		hpDrain /= 2
	}

	diff.HPMod = hpDrain
	diff.CircleRadius = DifficultyRate(cs, 54.4, 32, 9.6) * 1.00041 //some weird allowance osu has
	diff.Preempt = DifficultyRate(ar, 1800, 1200, 450)
	diff.FadeIn = DifficultyRate(ar, 1200, 800, 300)
	diff.Hit50 = int64(DifficultyRate(od, 200, 150, 100))
	diff.Hit100 = int64(DifficultyRate(od, 140, 100, 60))
	diff.Hit300 = int64(DifficultyRate(od, 80, 50, 20))
	diff.SpinnerRatio = DifficultyRate(od, 3, 5, 7.5)
}

func (diff *Difficulty) SetMods(mods Modifier) {
	diff.Mods = mods
	diff.calculate()
}

func (diff *Difficulty) CheckModActive(mods Modifier) bool {
	return diff.Mods&mods > 0
}

func (diff *Difficulty) GetModifiedTime(time float64) float64 {
	if diff.Mods&DoubleTime > 0 {
		return time / 1.5
	} else if diff.Mods&HalfTime > 0 {
		return time / 0.75
	} else {
		return time
	}
}

func (diff *Difficulty) GetHPDrain() float64 {
	return diff.hpDrain
}

func (diff *Difficulty) SetHPDrain(hpDrain float64) {
	diff.hpDrain = hpDrain
	diff.calculate()
}

func (diff *Difficulty) GetCS() float64 {
	return diff.cs
}

func (diff *Difficulty) SetCS(cs float64) {
	diff.cs = cs
	diff.calculate()
}

func (diff *Difficulty) GetOD() float64 {
	return diff.od
}

func (diff *Difficulty) SetOD(od float64) {
	diff.od = od
	diff.calculate()
}

func (diff *Difficulty) GetAR() float64 {
	return diff.ar
}

func (diff *Difficulty) SetAR(ar float64) {
	diff.ar = ar
	diff.calculate()
}

func DifficultyRate(diff, min, mid, max float64) float64 {
	diff = float64(float32(diff))

	if diff > 5 {
		return mid + (max-mid)*(diff-5)/5
	}
	if diff < 5 {
		return mid - (mid-min)*(5-diff)/5
	}
	return mid
}
