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

const GRID_WIDTH        = 200
const GRID_HEIGHT       = 150
const WIN_TITLE         = "Game of Life"
const CELL_COLOR_R      = 255
const CELL_COLOR_G      = 100
const CELL_COLOR_B      = 0
const BG_COLOR_R        = 0
const BG_COLOR_G        = 0
const BG_COLOR_B        = 0

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

    for y:=0; y < GRID_HEIGHT; y++ {
        for x:=0; x < GRID_WIDTH; x++ {
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
    cellToInt := func(x int, y int) int {
        // Return 0 if out of bounds
        if x < 0 || x >= GRID_WIDTH || y < 0 || y >= GRID_HEIGHT {
            return 0;
        }

        if mat[y][x].isAlive {
            return 1
        } else {
            return 0
        }
    }

    count := 0

    // Top
    count += cellToInt(pos.x-1, pos.y-1)
    count += cellToInt(pos.x-1, pos.y+0)
    count += cellToInt(pos.x-1, pos.y+1)

    // Middle
    count += cellToInt(pos.x+0, pos.y-1)
    count += cellToInt(pos.x+0, pos.y+1)

    // Bottom
    count += cellToInt(pos.x+1, pos.y-1)
    count += cellToInt(pos.x+1, pos.y+0)
    count += cellToInt(pos.x+1, pos.y+1)

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
    err = rend.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
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
                } else if key == sdl.K_DELETE {
                    matrix = Matrix{}
                    g_genCount = 1
                    isSimulating = false
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
                    pixels[index+1] = CELL_COLOR_B
                    pixels[index+2] = CELL_COLOR_G
                    pixels[index+3] = CELL_COLOR_R
                } else {
                    pixels[index+1] = BG_COLOR_B
                    pixels[index+2] = BG_COLOR_G
                    pixels[index+3] = BG_COLOR_R
                }
            }
        }
        matTex.Unlock()

        rend.Copy(matTex, nil, nil)

        { // Render preview cell under the cursor
            mx, my, _ := sdl.GetMouseState()
            cx := int(float32(mx)/cellW)
            cy := int(float32(my)/cellH)
            if cx < GRID_WIDTH && cy < GRID_HEIGHT {
                if matrix[cy][cx].isAlive {
                    rend.SetDrawColor(BG_COLOR_R, BG_COLOR_G, BG_COLOR_B, 100)
                } else {
                    rend.SetDrawColor(CELL_COLOR_R, CELL_COLOR_G, CELL_COLOR_B, 100)
                }
                x := int32(float32(int32(float32(mx)/cellW))*cellW)
                y := int32(float32(int32(float32(my)/cellH))*cellH)
                rend.FillRect(&sdl.Rect{X: x, Y: y, W: int32(cellW), H: int32(cellH)})
            }

        }

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
