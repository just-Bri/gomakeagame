// Package project maps simple 3D placeholder shapes into 2D screen space.
// Useful for isometric / soft-3D games before full g3d assets exist.
package project

import "math"

// Iso is a classic 2:1 isometric projector with an optional camera offset.
type Iso struct {
	// Origin is the screen-space position of world (0,0,0).
	OriginX, OriginY float64
	// Scale is pixels per world unit along an isometric axis (zoom).
	Scale float64
	// CamX/CamZ are world-space look targets subtracted before projection.
	CamX, CamZ float64
	// Yaw rotates the view around world +Y (radians). 0 is the classic SE-facing iso.
	Yaw float64
}

// DefaultIso returns a projector centered-ish for a 1280×720 virtual canvas.
func DefaultIso() Iso {
	return Iso{OriginX: 640, OriginY: 280, Scale: 28}
}

// WorldToScreen projects a world point (Y up) to screen pixels.
func (p Iso) WorldToScreen(x, y, z float64) (sx, sy float64) {
	x, z = p.rotateXZ(x-p.CamX, z-p.CamZ)
	sx = p.OriginX + (x-z)*p.Scale
	sy = p.OriginY + (x+z)*p.Scale*0.5 - y*p.Scale
	return sx, sy
}

// ScreenToGround casts a screen point onto the Y=0 plane.
func (p Iso) ScreenToGround(sx, sy float64) (x, z float64) {
	// Invert classic iso, then undo yaw:
	// sx = Ox + (x-z)*S
	// sy = Oy + (x+z)*S/2
	u := (sx - p.OriginX) / p.Scale
	v := (sy - p.OriginY) / (p.Scale * 0.5)
	rx := (u + v) / 2
	rz := (v - u) / 2
	x, z = p.unrotateXZ(rx, rz)
	return x + p.CamX, z + p.CamZ
}

// DepthKey sorts draw order (painter's algorithm): larger = further back.
func (p Iso) DepthKey(x, y, z float64) float64 {
	x, z = p.rotateXZ(x-p.CamX, z-p.CamZ)
	return x + z - y*0.01
}

// rotateXZ applies Yaw to a camera-relative XZ vector.
func (p Iso) rotateXZ(x, z float64) (rx, rz float64) {
	if p.Yaw == 0 {
		return x, z
	}
	c, s := math.Cos(p.Yaw), math.Sin(p.Yaw)
	return x*c - z*s, x*s + z*c
}

// unrotateXZ inverts rotateXZ.
func (p Iso) unrotateXZ(rx, rz float64) (x, z float64) {
	if p.Yaw == 0 {
		return rx, rz
	}
	c, s := math.Cos(p.Yaw), math.Sin(p.Yaw)
	return rx*c + rz*s, -rx*s + rz*c
}

// BoxCorners returns the 8 corners of an axis-aligned box centered at (x,y,z)
// with size (w,h,d). Y is up; box sits with its bottom at y - h/2.
func BoxCorners(x, y, z, w, h, d float64) [8][3]float64 {
	hw, hh, hd := w/2, h/2, d/2
	return [8][3]float64{
		{x - hw, y - hh, z - hd},
		{x + hw, y - hh, z - hd},
		{x + hw, y - hh, z + hd},
		{x - hw, y - hh, z + hd},
		{x - hw, y + hh, z - hd},
		{x + hw, y + hh, z - hd},
		{x + hw, y + hh, z + hd},
		{x - hw, y + hh, z + hd},
	}
}

// FaceIndices lists box faces as quads (corner indices) back-to-front-ish.
var FaceIndices = [6][4]int{
	{0, 1, 5, 4}, // -Z
	{3, 2, 6, 7}, // +Z
	{0, 3, 7, 4}, // -X
	{1, 2, 6, 5}, // +X
	{0, 1, 2, 3}, // -Y bottom
	{4, 5, 6, 7}, // +Y top
}

// FaceShade returns a crude directional shade factor 0.55..1.0 for a face index.
func FaceShade(face int) float64 {
	switch face {
	case 5: // top
		return 1.0
	case 3, 1: // lit sides
		return 0.82
	case 2, 0:
		return 0.68
	default: // bottom
		return 0.55
	}
}

// Clamp01 clamps v into [0,1].
func Clamp01(v float64) float64 {
	return math.Max(0, math.Min(1, v))
}
