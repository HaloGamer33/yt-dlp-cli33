package main

import (
   "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/container"
    //"fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2"
    "log"
    "os/exec"
    "io"
    "syscall"
)

func main() {
    app_yt_dlp := app.New()

    main_window := app_yt_dlp.NewWindow("Yt-dlp Downloader") 
    main_window.Resize(fyne.NewSize(500, 100))
    main_window.SetMaster()
    main_window.CenterOnScreen()

//  ╭──────────────────────────────────────────────────────────╮
//  │                     Creating Widgets                     │
//  ╰──────────────────────────────────────────────────────────╯
    link_input := widget.NewEntry()
    link_input.SetPlaceHolder("Video Link...")

    var resolution_string string
    resolutions := []string{"2160","1440","1080","720","360","240","144"}
    resolution_select := widget.NewSelect(resolutions, func(value string) {
        resolution_string = value
    })
    resolution_select.PlaceHolder = "Resolution..."

    var videoformat_string string
    formats := []string{"avi","flv","gif","mkv","mov","mp4","webm","aac","aiff","alac","flac","m4a","mka","mp3","ogg","opus","vorbis","wav"}
    format_select := widget.NewSelect(formats, func(value string) {
        videoformat_string = value
    })
    format_select.PlaceHolder = "Video Format..."

//  ╭──────────────────────────────────────────────────────────╮
//  │                   Submit Button Logic                    │
//  ╰──────────────────────────────────────────────────────────╯
    submit_button := widget.NewButton(
        "Download",
        func() {
            videolink := link_input.Text

            var resolution_cmd string
            if resolution_select.SelectedIndex() != -1 {
                resolution_cmd = "-f bestvideo[height=" + resolution_string + "]+bestaudio"
            }

            videoformat_cmd := []string{"",""}
            if format_select.SelectedIndex() != -1 {
                videoformat_cmd[0] = "--recode-video"
                videoformat_cmd[1] = videoformat_string
            }

            args := []string{
                videolink,
                resolution_cmd,
                videoformat_cmd[0],
                videoformat_cmd[1],
            }


            cmd := exec.Command("./dependencies/yt-dlp", args...)
            log.Println(cmd)
            cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

            stderr, _ := cmd.StderrPipe()
            stdout, _ := cmd.StdoutPipe()

            cmd.Start()

            // Showing Downloading Window
            downloading_window := app_yt_dlp.NewWindow("Downloading...")
            downloading_window.SetContent(
                container.NewGridWithRows(2,
                    widget.NewLabel("Downloading..."),
                    widget.NewProgressBarInfinite(),
                ),
            )
            downloading_window.CenterOnScreen()
            downloading_window.Show()

            // Printing to terminal the outs
            slurp, _ := io.ReadAll(stderr)
            log.Printf("%s\n", slurp)
            slurp_1, _ := io.ReadAll(stdout)
            log.Printf("%s\n", slurp_1)

            cmd.Wait()
            downloading_window.Close()
        },
    )

//  ╭──────────────────────────────────────────────────────────╮
//  │             Setting content and showing app              │
//  ╰──────────────────────────────────────────────────────────╯
    content := container.NewVBox(
        link_input,
        resolution_select,
        format_select,
        submit_button,
    )

    main_window.SetContent(content)
    main_window.ShowAndRun()
}
