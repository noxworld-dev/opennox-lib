package things

import (
	"fmt"
	"io"
	"strings"
)

type Draw interface {
	isDraw()
}

type BaseDraw struct {
	Img ImageRef `json:"base"`
}

func (BaseDraw) isDraw() {}

type StaticDraw struct {
	Img ImageRef `json:"static"`
}

func (StaticDraw) isDraw() {}

type WeaponDraw struct {
	Img ImageRef `json:"weapon"`
}

func (WeaponDraw) isDraw() {}

type ArmorDraw struct {
	Img ImageRef `json:"armor"`
}

func (ArmorDraw) isDraw() {}

type StaticRandomDraw struct {
	Imgs []ImageRef `json:"random"`
}

func (StaticRandomDraw) isDraw() {}

type DoorDraw struct {
	Imgs []ImageRef `json:"door"`
}

func (DoorDraw) isDraw() {}

type AnimateDraw struct {
	Anim Animation `json:"anim"`
}

func (AnimateDraw) isDraw() {}

type GlyphDraw struct {
	Anim Animation `json:"glyph"`
}

func (GlyphDraw) isDraw() {}

type WeaponAnimateDraw struct {
	Anim Animation `json:"weapon_anim"`
}

func (WeaponAnimateDraw) isDraw() {}

type ArmorAnimateDraw struct {
	Anim Animation `json:"armor_anim"`
}

func (ArmorAnimateDraw) isDraw() {}

type FlagDraw struct {
	Anim Animation `json:"flag"`
}

func (FlagDraw) isDraw() {}

type SphericalShieldDraw struct {
	Anim Animation `json:"spherical_shield"`
}

func (SphericalShieldDraw) isDraw() {}

type SummonEffectDraw struct {
	Anim Animation `json:"summon"`
}

func (SummonEffectDraw) isDraw() {}

type ConditionalAnimateDraw struct {
	Anims []Animation `json:"cond_anim"`
}

func (ConditionalAnimateDraw) isDraw() {}

type MonsterAnimationType byte

type MonsterAnimation struct {
	Type         MonsterAnimationType `json:"type"`
	Sound        string               `json:"sound,omitempty"`
	Field8       string               `json:"field_8,omitempty"`
	FramesPerDir byte                 `json:"frames_per_dir"`
	Field10      byte                 `json:"field_10,omitempty"`
	Kind         AnimationKind        `json:"kind"`
	Frames       []ImageRef           `json:"frames"`
}

type MonsterDraw struct {
	Anims []MonsterAnimation `json:"monster"`
}

func (MonsterDraw) isDraw() {}

type MaidenDraw struct {
	Anims []MonsterAnimation `json:"maiden"`
}

func (MaidenDraw) isDraw() {}

type MonsterGeneratorDraw struct {
	Anims []Animation `json:"monster_gen"`
}

func (MonsterGeneratorDraw) isDraw() {}

type PlayerAnim struct {
	FramesPerDir byte   `json:"frames_per_dir"`
	Field8       byte   `json:"field_8,omitempty"`
	Field12      string `json:"field_12,omitempty"`
}

type PlayerAnimType string
type PlayerAnimPart string

type PlayerDraw struct {
	Anims map[PlayerAnimType]PlayerAnim `json:"player_anim"`
	Parts map[PlayerAnimPart][]ImageRef `json:"player_part"`
}

func (PlayerDraw) isDraw() {}

func (f *Reader) readAnimation() (*Animation, error) {
	frames, err := f.readU8()
	if err != nil {
		return nil, err
	}
	ani := &Animation{Frames: make([]ImageRef, 0, frames)}
	v2, err := f.readU8()
	if err != nil {
		return nil, err
	}
	ani.Field = v2
	loop, err := f.readBytes8()
	if err != nil {
		return nil, err
	}
	if err := ani.Kind.UnmarshalText(loop); err != nil {
		return nil, err
	}
	for i := 0; i < int(frames); i++ {
		fr, err := f.readImageRef()
		if err != nil {
			return nil, err
		}
		ani.Frames = append(ani.Frames, *fr)
	}
	return ani, nil
}

func (f *Reader) skipAnimation() error {
	n, err := f.readU8()
	if err != nil {
		return err
	}
	if err := f.skip(1); err != nil {
		return err
	}
	if err := f.skipBytes8(); err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		if err := f.skipImageRef(); err != nil {
			return err
		}
	}
	return err
}

func (f *Reader) readAnimations8() ([]Animation, error) {
	n, err := f.readU8()
	if err != nil {
		return nil, err
	}
	out := make([]Animation, 0, n)
	for i := 0; i < int(n); i++ {
		ani, err := f.readAnimation()
		if err != nil {
			return nil, err
		}
		out = append(out, *ani)
	}
	return out, nil
}

func (f *Reader) skipAnimations8() error {
	n, err := f.readU8()
	if err != nil {
		return err
	}
	for i := 0; i < int(n); i++ {
		if err := f.skipAnimation(); err != nil {
			return err
		}
	}
	return nil
}

func (f *Reader) readMonsterDraw() (MonsterDraw, error) {
	d := MonsterDraw{}
	for {
		sect, err := f.readSect()
		if err == io.EOF {
			return d, io.ErrUnexpectedEOF
		} else if err != nil {
			return d, err
		}
		switch sect {
		default:
			return d, fmt.Errorf("unsupported monster draw sect: %q", sect)
		case "END ":
			return d, nil
		case "STAT":
			typ, err := f.readU8()
			if err != nil {
				return d, err
			}
			snd, err := f.readString8()
			if err != nil {
				return d, err
			}
			fld8, err := f.readString8()
			if err != nil {
				return d, err
			}
			framesN, err := f.readU8()
			if err != nil {
				return d, err
			}
			fld10, err := f.readU8()
			if err != nil {
				return d, err
			}
			kind, err := f.readBytes8()
			if err != nil {
				return d, err
			}
			ani := MonsterAnimation{
				Type:         MonsterAnimationType(typ),
				Sound:        snd,
				Field8:       strings.TrimRight(fld8, "\x00"),
				FramesPerDir: framesN,
				Field10:      fld10,
				Frames:       make([]ImageRef, 0, 8*int(framesN)),
			}
			if err = ani.Kind.UnmarshalText(kind); err != nil {
				return d, err
			}
			for i := 0; i < 8; i++ {
				for j := 0; j < int(framesN); j++ {
					ref, err := f.readImageRef()
					if err != nil {
						return d, err
					}
					ani.Frames = append(ani.Frames, *ref)
				}
			}
			d.Anims = append(d.Anims, ani)
		}
	}
}

func (f *Reader) readPlayerDraw() (PlayerDraw, error) {
	d := PlayerDraw{
		Anims: make(map[PlayerAnimType]PlayerAnim),
		Parts: make(map[PlayerAnimPart][]ImageRef),
	}
	var lastFrames int
	for {
		sect, err := f.readSect()
		if err == io.EOF {
			return d, io.ErrUnexpectedEOF
		} else if err != nil {
			return d, err
		}
		switch sect {
		default:
			return d, fmt.Errorf("unsupported monster draw sect: %q", sect)
		case "END ":
			return d, nil
		case "STAT":
			name, err := f.readString8()
			if err != nil {
				return d, err
			}
			framesN, err := f.readU8()
			if err != nil {
				return d, err
			}
			lastFrames = int(framesN)
			fld8, err := f.readU8()
			if err != nil {
				return d, err
			}
			fld12, err := f.readString8()
			if err != nil {
				return d, err
			}
			ani := PlayerAnim{
				FramesPerDir: framesN,
				Field8:       fld8,
				Field12:      fld12,
			}
			d.Anims[PlayerAnimType(name)] = ani
		case "SEQU":
			name, err := f.readString8()
			if err != nil {
				return d, err
			}
			frames := make([]ImageRef, 0, 8*lastFrames)
			for i := 0; i < 8; i++ {
				for j := 0; j < lastFrames; j++ {
					ref, err := f.readImageRef()
					if err != nil {
						return d, err
					}
					frames = append(frames, *ref)
				}
			}
			d.Parts[PlayerAnimPart(name)] = frames
		}
	}
}

func (f *Reader) skipThingDraw() error {
	dname, err := f.readString8()
	if err != nil {
		return err
	}
	sectSz, err := f.readU64align()
	if err != nil {
		return err
	}
	switch dname {
	case "StaticDraw", "WeaponDraw", "ArmorDraw", "BaseDraw":
		if err := f.skipImageRef(); err != nil {
			return err
		}
	case "StaticRandomDraw", "DoorDraw":
		if err := f.skipImageRefs8(); err != nil {
			return err
		}
	case "AnimateDraw", "GlyphDraw", "WeaponAnimateDraw", "ArmorAnimateDraw",
		"FlagDraw", "SphericalShieldDraw", "SummonEffectDraw":
		if err := f.skipAnimation(); err != nil {
			return err
		}
	case "ConditionalAnimateDraw", "MonsterGeneratorDraw":
		if err := f.skipAnimations8(); err != nil {
			return err
		}
	case "AnimateStateDraw":
		if err := f.skipThingAnimStateDraw(); err != nil {
			return err
		}
	case "VectorAnimateDraw", "ReleasedSoulDraw":
		if err := f.skipThingAnimVectorDraw(); err != nil {
			return err
		}
	case "MonsterDraw", "MaidenDraw":
		if err := f.skipThingMonsterDraw(); err != nil {
			return err
		}
	case "PlayerDraw":
		if err := f.skipThingPlayerDraw(); err != nil {
			return err
		}
	case "SlaveDraw", "BoulderDraw", "ArrowDraw", "WeakArrowDraw", "HarpoonDraw":
		if err := f.skipImageRefs8(); err != nil {
			return err
		}
	default:
		if err := f.skip(int(sectSz)); err != nil {
			return err
		}
	}
	return nil
}

func (f *Reader) readThingDraw() (Draw, error) {
	dname, err := f.readString8()
	if err != nil {
		return nil, err
	}
	sectSz, err := f.readU64align()
	if err != nil {
		return nil, err
	}
	switch dname {
	case "StaticDraw", "WeaponDraw", "ArmorDraw", "BaseDraw":
		img, err := f.readImageRef()
		if err != nil {
			return nil, err
		}
		switch dname {
		case "StaticDraw":
			return StaticDraw{Img: *img}, nil
		case "WeaponDraw":
			return WeaponDraw{Img: *img}, nil
		case "ArmorDraw":
			return ArmorDraw{Img: *img}, nil
		case "BaseDraw":
			return BaseDraw{Img: *img}, nil
		}
	case "StaticRandomDraw", "DoorDraw":
		imgs, err := f.readImageRefs8()
		if err != nil {
			return nil, err
		}
		switch dname {
		case "StaticRandomDraw":
			return StaticRandomDraw{Imgs: imgs}, nil
		case "DoorDraw":
			return DoorDraw{Imgs: imgs}, nil
		}
	case "AnimateDraw", "GlyphDraw", "WeaponAnimateDraw", "ArmorAnimateDraw",
		"FlagDraw", "SphericalShieldDraw", "SummonEffectDraw":
		anim, err := f.readAnimation()
		if err != nil {
			return nil, err
		}
		switch dname {
		case "AnimateDraw":
			return AnimateDraw{Anim: *anim}, nil
		case "GlyphDraw":
			return GlyphDraw{Anim: *anim}, nil
		case "WeaponAnimateDraw":
			return WeaponAnimateDraw{Anim: *anim}, nil
		case "ArmorAnimateDraw":
			return ArmorAnimateDraw{Anim: *anim}, nil
		case "FlagDraw":
			return FlagDraw{Anim: *anim}, nil
		case "SphericalShieldDraw":
			return SphericalShieldDraw{Anim: *anim}, nil
		case "SummonEffectDraw":
			return SummonEffectDraw{Anim: *anim}, nil
		}
	case "ConditionalAnimateDraw", "MonsterGeneratorDraw":
		anis, err := f.readAnimations8()
		if err != nil {
			return nil, err
		}
		switch dname {
		case "ConditionalAnimateDraw":
			return ConditionalAnimateDraw{Anims: anis}, nil
		case "MonsterGeneratorDraw":
			return MonsterGeneratorDraw{Anims: anis}, nil
		}
	case "AnimateStateDraw":
		if err := f.skipThingAnimStateDraw(); err != nil {
			return nil, err // FIXME
		}
	case "VectorAnimateDraw", "ReleasedSoulDraw":
		if err := f.skipThingAnimVectorDraw(); err != nil {
			return nil, err // FIXME
		}
	case "MonsterDraw", "MaidenDraw":
		d, err := f.readMonsterDraw()
		if err != nil {
			return nil, err
		}
		switch dname {
		case "MonsterDraw":
			return d, nil
		case "MaidenDraw":
			return MaidenDraw(d), nil
		}
	case "PlayerDraw":
		d, err := f.readPlayerDraw()
		if err != nil {
			return nil, err
		}
		return d, nil
	case "SlaveDraw", "BoulderDraw", "ArrowDraw", "WeakArrowDraw", "HarpoonDraw":
		if err := f.skipImageRefs8(); err != nil {
			return nil, err // FIXME
		}
	default:
		if err := f.skip(int(sectSz)); err != nil {
			return nil, err
		}
	}
	return nil, nil
}
