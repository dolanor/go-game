// Package mathlib has golang implementations of matrix stuff
// needed for 3d games
package mathlib

import "math"

// PerspectiveMat returns a 4x4 projection matrix created using
// the parameters fovy, aspect, near, far.
//Usage: PerspectiveMat(math.PI / 2, WIDTH/HEIGHT,0.1,100)
func PerspectiveMat(fovy, aspect, near, far float64) [4][4]float64 {
	var m [4][4]float64
	tanHalfFovy := math.Tan(0.5 * fovy)
	m[0][0] = 1 / (aspect * tanHalfFovy)
	m[1][1] = 1 / (tanHalfFovy)
	m[2][2] = +(far + near) / (far - near)
	m[2][3] = +1
	m[3][2] = -2 * far * near / (far - near)
	return m
}

// RotationMat returns a 3x3 rotation matrix with parameters
// angle and rotation origin. When a position is multiplied
// by it, the result is the original position rotated in 3-dimensions.
func RotationMat(ang float64, axis [3]float64) (rot [3][3]float64) {
	c := math.Cos(ang)
	s := math.Sin(ang)

	alen := math.Sqrt(axis[0]*axis[0] + axis[1]*axis[1] + axis[2]*axis[2])
	a := [3]float64{axis[0] / alen, axis[1] / alen, axis[2] / alen}
	t := [3]float64{a[0] * (1 - c), a[1] * (1 - c), a[2] * (1 - c)}

	rot[0][0] = c + t[0]*a[0]
	rot[0][1] = 0 + t[0]*a[1] + s*a[2]
	rot[0][2] = 0 + t[0]*a[2] - s*a[1]

	rot[1][0] = 0 + t[1]*a[0] - s*a[2]
	rot[1][1] = c + t[1]*a[1]
	rot[1][2] = 0 + t[1]*a[2] + s*a[0]

	rot[2][0] = 0 + t[2]*a[0] + s*a[1]
	rot[2][1] = 0 + t[2]*a[1] - s*a[0]
	rot[2][2] = c + t[2]*a[2]
	return
}

// MultiplyMat4414 does [4][4] x [1][4] -> [1][4]
// The operation is [J][I]  x [K][J] -> [K][I]
func MultiplyMat4414(a [4][4]float64, b [1][4]float64) (c [1][4]float64) {
	K, J, I := 1, 4, 4
	for k := 0; k < K; k++ {
		for j := 0; j < J; j++ {
			for i := 0; i < I; i++ {
				c[k][i] += a[j][i] * b[k][j]
			}
		}
	}
	return
}

// MultiplyMatByVec4 does [4][4] x vec4 -> vec4
func MultiplyMatByVec4(a [4][4]float64, b [4]float64) (c [4]float64) {
	J, I := 4, 4
	for j := 0; j < J; j++ {
		for i := 0; i < I; i++ {
			c[i] += a[j][i] * b[j]
		}
	}
	return
}

// MultiplyMat3313 does [3][3] x [1][3] -> [1][3]
// The operation is [J][I]  x [K][J] -> [K][I]
func MultiplyMat3313(a [3][3]float64, b [1][3]float64) (c [1][3]float64) {
	K, J, I := 1, 3, 3
	for k := 0; k < K; k++ {
		for j := 0; j < J; j++ {
			for i := 0; i < I; i++ {
				c[k][i] += a[j][i] * b[k][j]
			}
		}
	}
	return
}

// MultiplyMatByVec4 does [4][4] x vec4 -> vec4
func MultiplyMatVec3(a [3][3]float64, b [3]float64) (c [3]float64) {
	J, I := 3, 3
	for j := 0; j < J; j++ {
		for i := 0; i < I; i++ {
			c[i] += a[j][i] * b[j]
		}
	}
	return
}
