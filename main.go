package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type rules struct {
	StopAtNextStation bool
	Iteration int
	WithinStationBoundary bool
	Failed bool
	Score int 
	Honked bool
	ShowInstructions bool
}

const screenWidth = int32(1080)
const screenHeight = int32(720)
var raindropAvgHeight = int32(5)
var speed = float32(0.0)

var rule = &rules{
	StopAtNextStation: false,
	Iteration: 0,
	WithinStationBoundary:false,
	Failed: false,
	Score: 0,
	Honked:false,
	ShowInstructions:true,
}

func main() {

	// ------------------ Initialize Audio ---------------------

	rl.InitAudioDevice()

    rainMusic :=rl.LoadMusicStream("./assets/audio/rn.mp3")
    rl.PlayMusicStream(rainMusic)
    defer rl.UnloadMusicStream(rainMusic)

    lightning :=rl.LoadSound("./assets/audio/lightning.mp3")
    defer rl.UnloadSound(lightning)

    trainMusic :=rl.LoadMusicStream("./assets/audio/train.mp3")
    rl.PlayMusicStream(trainMusic)
    rl.SetMusicVolume(trainMusic, 2)
    defer rl.UnloadMusicStream(trainMusic)

    horn :=rl.LoadMusicStream("./assets/audio/horn.mp3")
    rl.PlayMusicStream(horn)
    defer rl.UnloadMusicStream(horn)

	// ----------- Initialize Window -------------
	rl.InitWindow(screenWidth, screenHeight, "3d camera first person rain")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	camera := rl.NewCamera3D(
		rl.NewVector3(20,4,4),
		rl.NewVector3(20,-1.8,10000000),
		rl.NewVector3(0,1,0),
		60.0,
		rl.CameraPerspective,
	)

	// -------------- Load and Store 3D Models ------------------

    red_signal := rl.LoadModel("./assets/models/red_signal.glb")
    defer rl.UnloadModel(red_signal)
    green_signal := rl.LoadModel("./assets/models/green_signal.glb")
    defer rl.UnloadModel(green_signal)
    terrain := rl.LoadModel("./assets/models/terrain.glb")
    defer rl.UnloadModel(terrain)
    tree := rl.LoadModel("./assets/models/tree.glb")
    defer rl.UnloadModel(tree)
    track := rl.LoadModel("./assets/models/track.glb")
    defer rl.UnloadModel(track)
    train_station := rl.LoadModel("./assets/models/train_station.glb")
    defer rl.UnloadModel(train_station)
    sign := rl.LoadModel("./assets/models/sign.glb")
    defer rl.UnloadModel(sign)
    electricity := rl.LoadModel("./assets/models/electricity.glb")
    defer rl.UnloadModel(electricity)
    track_bent := rl.LoadModel("./assets/models/track_bent.glb")
    defer rl.UnloadModel(track_bent)
    track_bent_r := rl.LoadModel("./assets/models/track_bent_r.glb")
    defer rl.UnloadModel(track_bent_r)
    mountains := rl.LoadModel("./assets/models/mountains.glb")
    defer rl.UnloadModel(mountains)


	models := make(map[string]rl.Model)
	models["red_signal"]=red_signal
	models["green_signal"]=green_signal
	models["terrain"]=terrain
	models["tree"]=tree
	models["track"]=track
	models["train_station"]=train_station
	models["sign"]=sign
	models["electricity"]=electricity
	models["track_bent"]=track_bent
	models["track_bent_r"]=track_bent_r
	models["mountains"]=mountains

	clearModels := func(m map[string]rl.Model) {
        for key := range m {
            delete(m, key)
        }
    }
    defer clearModels(models)

	// --------------- Start Game -------------------------------

	for !rl.WindowShouldClose() {
	 	updateCamera(&camera)
		rainSound(lightning,rainMusic)
		drawGame(&camera,trainMusic,models)
		showInstructions(camera)
		controls(horn)
        }
}

func updateCamera(camera *rl.Camera3D) {
	_ = camera
	//rl.UpdateCamera(camera,rl.CameraCustom)  if you want to change the perspentive
	rl.BeginDrawing()
	defer rl.EndDrawing()
}

func controls(horn rl.Music) {

	if rl.IsKeyDown(rl.KeyH) {
		rl.UpdateMusicStream(horn)				
	}

	if rl.IsKeyReleased(rl.KeyH) {
		rl.SeekMusicStream(horn,float32(0))
	}

	if rl.IsKeyDown(rl.KeyW) { 
		speed += 0.005
	}

	if rl.IsKeyDown(rl.KeyB) {
		if speed > 0 {
			speed -= 0.02
		}
	}

	if rl.IsMouseButtonDown(rl.MouseRightButton) {
	   speed = 0
	}

	if rl.IsKeyDown(rl.KeyS){
		speed -= 0.005
	}

	if (rl.IsKeyDown(rl.KeyW) && speed > 5.1) {
		rl.DrawText("Max Speed", screenWidth / 2 - 100, screenHeight / 2 - 100, 20, rl.Red)
	}

	if (rl.IsKeyDown(rl.KeyH) && rule.WithinStationBoundary && !rule.Honked) {
		rule.Honked = true
		rule.Score += 100
	}


}

func drawGame(camera *rl.Camera3D,trainMusic rl.Music, models map[string]rl.Model) {
	rl.BeginMode3D(*camera)
	defer rl.EndMode3D()

	camera.Position.Z += float32(speed)

	if speed > 0 {
		speed -= 0.001
		rl.UpdateMusicStream(trainMusic)
		if speed < 1 {
			rl.SetMusicPitch(trainMusic,speed)	
		}
	}

	if (camera.Position.Z > 2000) {
		camera.Position.Z = 0
		rule.Iteration++
	}

	raindropAvgHeight--
	if raindropAvgHeight < 0 {
		raindropAvgHeight = 5
	}

	for rain_x := int32(-10); rain_x < 10; rain_x++ {
	  for rain_y := int32(-10); rain_y < 10; rain_y++ {
		rl.DrawCube(
			rl.NewVector3(
				camera.Position.X + 0.5 + float32(rain_x + rl.GetRandomValue(1, 2)),
				float32(raindropAvgHeight + rl.GetRandomValue(0, 3)),
				camera.Position.Z + float32(rain_y + rl.GetRandomValue(0, 2)),
			),
			0.01,
			0.2,
			0.01,
			rl.Blue,
		)
	  }	
	}

	for x := -4;x<4; x++ {
		for z := -1; z <100; z++ {
			rl.DrawModel(models["tree"],rl.NewVector3(float32(x)*40,12,float32(z)*40),0.3,rl.Brown)
		}
	}

	for z := 0 ; z<40 ; z++ {
		rl.DrawModel(models["mountains"],rl.NewVector3(-300,30,float32(z)*200),2,rl.LightGray)
		rl.DrawModel(models["mountains"],rl.NewVector3(150,30,float32(z)*200),2,rl.LightGray)
	}

	for z := 0 ; z<500 ; z++ {
		rl.DrawModel(models["track"],rl.NewVector3(20,1.5,float32(z)*17),0.15,rl.DarkGray)
		rl.DrawModel(models["track"],rl.NewVector3(10,1.5,float32(z)*17),0.15,rl.DarkGray)
	}

	for z := 1 ; z<50 ; z++ {
		rl.DrawModel(models["electricity"],rl.NewVector3(23.5,6,float32(z)*50),0.1,rl.Gray)
	}

	rl.DrawModel(models["sign"],rl.NewVector3(30,2,200),0.1,rl.Gray)
	for z := 1 ; z<50 ; z++ {
		rl.DrawModel(models["sign"],rl.NewVector3(30,2,float32(z)*550),0.1,rl.Gray)
	}

	rl.DrawModel(models["train_station"], rl.NewVector3(35, 3.4, -9), 0.3, rl.Gray)
	rl.DrawModel(models["train_station"], rl.NewVector3(35, 3.4, 1990), 0.3, rl.Gray)

	if (rule.StopAtNextStation) {
		rl.DrawModel(models["red_signal"], rl.NewVector3(28, 2.4, 2040), 0.1, rl.Gray)
		rl.DrawModel(models["red_signal"], rl.NewVector3(28, 2.4, 40), 0.1, rl.Gray)
	} else {
		rl.DrawModel(models["green_signal"], rl.NewVector3(28, 2.4, 2040), 0.1, rl.Gray)
		rl.DrawModel(models["green_signal"], rl.NewVector3(28, 2.4, 40), 0.1, rl.Gray)
	}

	rl.DrawCube(rl.NewVector3(11, 0.1, 0.0), 6, 0.01, 7000, rl.Gray)
	rl.DrawModel(models["track_bent_r"], rl.NewVector3(15.7, 1.5, 50), 0.15, rl.Gray)
	rl.DrawModel(models["track_bent"], rl.NewVector3(15.7, 1.5, 1720), 0.15, rl.Gray)
	rl.DrawCube(rl.NewVector3(20.8, 0.1, 0.0), 6, 0.01, 7000, rl.Gray)
	rl.DrawCube(rl.NewVector3(21, 8, 0.0), 0.1, 0.1, 7000, rl.Black)
	rl.DrawCube(rl.NewVector3(0.0, 0, 0.0), 500, 0.01, 7000, rl.DarkBrown)

    if (camera.Position.Z > 1967 || camera.Position.Z < 7) {
		rule.WithinStationBoundary = true
	} else {
		rule.WithinStationBoundary = false
	}

	if (rule.Iteration > 0 && math.Remainder(float64(rule.Iteration), 2) == 0) {
		rule.StopAtNextStation = true
	}
	if (rule.StopAtNextStation && rule.WithinStationBoundary && speed <= 0) {
		rule.StopAtNextStation = false // after 1 station
		rule.Iteration = 0
		// WIN
	}
	if (rule.StopAtNextStation && camera.Position.Z > 1990 && speed != 0) {
		rule.Failed = true
	}

	if (!rule.WithinStationBoundary) {
		rule.Honked = false
	}

}

func rainSound(lightning rl.Sound,rainMusic rl.Music) {

	if (rl.GetRandomValue(0, 600) == 30) {
		rl.ClearBackground(rl.RayWhite)
		rl.PlaySound(lightning)
	} else {
		rl.ClearBackground(rl.Gray)
	}

	rl.UpdateMusicStream(rainMusic)
}

func showInstructions(camera rl.Camera3D) {
	if (rule.ShowInstructions) {
		rl.DrawRectangle(10,10,250,110,rl.Fade(rl.SkyBlue,0.5))
		rl.DrawRectangleLines(10,10,250,110,rl.Blue)
		rl.DrawText("Train controls:", 20, 20, 10, rl.Black)
		rl.DrawText("- Go forward: W, Go back : S", 40, 40, 10, rl.DarkGray)
		rl.DrawText("- Press H to Honk, Press B for breaks", 40, 60, 10, rl.DarkGray)
		rl.DrawText("- Honk at stations to increase score", 40, 80, 10, rl.DarkGray)
		rl.DrawText("- Press I to hide these instructions", 40, 100, 10, rl.DarkGray)
	}

	if (rule.Failed) {
		rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.Fade(rl.DarkGray,0.5))
		rl.DrawRectangleLines(10, 10, 250, 70, rl.DarkGray)
		rl.DrawText("Failed", screenWidth / 2 - 150, screenHeight / 2 - 100, 200, rl.Red)
	}
	rl.DrawText(fmt.Sprintf("Score: %d",rule.Score),  screenWidth - 220, 10, 20, rl.Black)
	rl.DrawText(fmt.Sprintf("Speed: %d Km/h", int(speed * 20)),screenWidth - 220, 30, 20, rl.Black)
	rl.DrawText(fmt.Sprintf("Next station: %d m",int(2000 - camera.Position.Z)),  screenWidth - 220, 50, 20, rl.Black)

	if (rule.StopAtNextStation) {
		rl.DrawText("Status: Stop at", screenWidth-220, 70, 20, rl.Red)
		rl.DrawText("          next station", screenWidth-220, 90, 20, rl.Red)
		rl.DrawText("Stop at Next station", screenWidth-500, 300, 25, rl.Red)
	} else {
		rl.DrawText("Status: Don't stop", screenWidth-220, 70, 20, rl.Green)
	}

	if (camera.Position.Z < 0) {
		camera.Position.Z = 0.1
		speed = 0
		rl.DrawText("Wrong Direction", screenWidth / 2 - 100, screenHeight / 2 - 100, 20, rl.Red)
	}

	if rl.IsKeyPressed(rl.KeyI) {
		if rule.ShowInstructions {
			rule.ShowInstructions = false
		} else {
			rule.ShowInstructions = true
		}
	}
}