package main

import (
    "fmt"
    "github.com/veandco/go-sdl2/sdl"
)

func CHECK_ERR(err error) {
    if err != nil {
        panic(err)
    }
}

const GRID_WIDTH = 200
const GRID_HEIGHT = 150
const WIN_TITLE = "Game of Life"

var g_genCount = 1

type Matrix [GRID_HEIGHT][GRID_WIDTH]MatrixCell;

type MatrixPos struct {
    x int
    y int
}

type MatrixCell struct {
    isAlive bool
}

func SimGeneration(mat *Matrix) Matrix {
    // Do the simulation in a copy
    mat2 := *mat

    for y:=1; y < GRID_HEIGHT-1; y++ {
        for x:=1; x < GRID_WIDTH-1; x++ {
            neighCnt := CountNeighb(mat, &MatrixPos{x: x, y: y})

            if neighCnt < 2 {
                mat2[y][x].isAlive = false
            } else if neighCnt == 2 {
                // Nothing happens
            } else if neighCnt == 3 {
                mat2[y][x].isAlive = true
            } else {
                mat2[y][x].isAlive = false
            }
        }
    }

    g_genCount++

    return mat2
}

func CountNeighb(mat *Matrix, pos *MatrixPos) int {
    cellToInt := func(cell *MatrixCell) int {
        if cell.isAlive {
            return 1
        } else {
            return 0
        }
    }

    count := 0

    // Top
    count += cellToInt(&(*mat)[pos.y-1][pos.x-1])
    count += cellToInt(&(*mat)[pos.y-1][pos.x+0])
    count += cellToInt(&(*mat)[pos.y-1][pos.x+1])

    // Middle
    count += cellToInt(&(*mat)[pos.y+0][pos.x-1])
    count += cellToInt(&(*mat)[pos.y+0][pos.x+1])

    // Bottom
    count += cellToInt(&(*mat)[pos.y+1][pos.x-1])
    count += cellToInt(&(*mat)[pos.y+1][pos.x+0])
    count += cellToInt(&(*mat)[pos.y+1][pos.x+1])

    return count
}

func main() {
    err := sdl.Init(sdl.INIT_VIDEO)
    CHECK_ERR(err)

    winRatio := float32(GRID_WIDTH)/GRID_HEIGHT
    winW := int32(1000*winRatio)
    winH := int32(1000)
    cellW := float32(winW)/GRID_WIDTH
    cellH := float32(winH)/GRID_HEIGHT
    win, err := sdl.CreateWindow(WIN_TITLE+" - Paused", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winW, winH, 0)
    CHECK_ERR(err)

    rend, err := sdl.CreateRenderer(win, 0, 0)
    CHECK_ERR(err)

    matrix := Matrix{}
    matTex, err := rend.CreateTexture(sdl.PIXELFORMAT_RGBX8888, sdl.TEXTUREACCESS_STREAMING, GRID_WIDTH, GRID_HEIGHT)
    CHECK_ERR(err)

    running := true
    isSimulating := false
    for running {
        var event sdl.Event
        for {
            event = sdl.PollEvent()
            if event == nil {
                break
            }
            switch event.GetType() {
            case sdl.QUIT:
                running = false

            case sdl.MOUSEBUTTONDOWN:
                x := int(float32(event.(*sdl.MouseButtonEvent).X)/cellW)
                y := int(float32(event.(*sdl.MouseButtonEvent).Y)/cellH)
                if x < GRID_WIDTH && y < GRID_HEIGHT {
                    matrix[y][x].isAlive = !matrix[y][x].isAlive
                }

            case sdl.KEYDOWN:
                key := event.(*sdl.KeyboardEvent).Keysym.Sym
                if key == sdl.K_SPACE {
                    isSimulating = !isSimulating
                } else if key == sdl.K_RETURN {
                    matrix = SimGeneration(&matrix)
                }
            }
        }
        if !running {
            break
        }

        if isSimulating {
            matrix = SimGeneration(&matrix)
        }

        pixels, _, err := matTex.Lock(nil) 
        CHECK_ERR(err)
        for y:=0; y < GRID_HEIGHT; y++ {
            for x:=0; x < GRID_WIDTH; x++ {
                index := (y*GRID_WIDTH+x)*4
                if matrix[y][x].isAlive {
                    pixels[index+1] = 0
                    pixels[index+2] = 255
                    pixels[index+3] = 0
                } else {
                    pixels[index+1] = 50
                    pixels[index+2] = 50
                    pixels[index+3] = 50
                }
            }
        }
        matTex.Unlock()

        rend.Copy(matTex, nil, nil)

        if isSimulating {
            win.SetTitle(fmt.Sprintf("%s - Simulating | Generation: %d", WIN_TITLE, g_genCount))
        } else {
            win.SetTitle(fmt.Sprintf("%s - Paused | Generation: %d", WIN_TITLE, g_genCount))
        }

        rend.Present()
        sdl.Delay(16)
    }

    err = rend.Destroy()
    CHECK_ERR(err)
    err = win.Destroy()
    CHECK_ERR(err)
    sdl.Quit()
}
