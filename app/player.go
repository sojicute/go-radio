package main

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v7/pkg/app"

	//"github.com/maxence-charriere/go-app/v7/pkg/errors"
	"strconv"
)

var (
	youtubeAPIReady bool
)

type player struct {
	app.Compo

	State playerState

	releaseIframe            func()
	releaseOnPlayerReady     func()
	releasePlayerStateChange func()

	onYouTubeIframeAPIReady bool
	youtube                 app.Value
	playing                 bool
	ready                   bool
}

func (p *player) setupYoutubePlayer() {
	onPlayerReady := app.FuncOf(p.onPlayerReady)
	p.releaseOnPlayerReady = onPlayerReady.Release

	onPlayerStateChange := app.FuncOf(p.onPlayerStateChange)
	p.releasePlayerStateChange = onPlayerStateChange.Release

	p.youtube = app.Window().
		Get("YT").
		Get("Player").
		New("youtube-player", map[string]interface{}{
			"events": map[string]interface{}{
				"onReady":       onPlayerReady,
				"onStateChange": onPlayerStateChange,
			},
		})
}

func (p *player) onPlayerReady(this app.Value, args []app.Value) interface{} {
	app.Dispatch(func() {
		p.setVolume(p.State.Volume)
		p.play()
	})
	return nil
}

func (p *player) onPlayerStateChange(this app.Value, args []app.Value) interface{} {
	playing := false
	ready := false

	switch args[0].Get("data").Int() {
	case 1:
		playing = true
		ready = true

	case -1, 0, 2:
		playing = false
		ready = true

	case 3:
		ready = false
	}
	app.Dispatch(func() {
		p.ready = ready
		p.playing = playing
		p.Update()
	})

	return nil
}

func (p *player) OnMount(ctx app.Context) {
	p.State.load()
	p.Update()

	if !youtubeAPIReady {
		onYouTubeIframeAPIReady := app.FuncOf(p.onYoutubeIframeAPIReady)
		p.releaseIframe = onYouTubeIframeAPIReady.Release
		app.Window().Set("onYouTubeIframeAPIReady", onYouTubeIframeAPIReady)
		return
	}

	app.Dispatch(p.setupYoutubePlayer)
}

func (p *player) OnDismount() {
	if p.releaseIframe != nil {
		p.releaseIframe()
	}

	if p.releaseOnPlayerReady != nil {
		p.releaseOnPlayerReady()
	}

	if p.releasePlayerStateChange != nil {
		p.releasePlayerStateChange()
	}
}

func (p *player) onYoutubeIframeAPIReady(this app.Value, args []app.Value) interface{} {
	app.Dispatch(func() {
		youtubeAPIReady = true
		p.setupYoutubePlayer()
	})

	return nil
}

func (p player) play() {
	p.youtube.Call("playVideo")
}

func (p player) pause() {
	p.youtube.Call("pauseVideo")
}

func (p *player) onPlay(ctx app.Context, e app.Event) {
	p.play()
}

func (p *player) onPause(ctx app.Context, e app.Event) {
	p.pause()
}

func (p *player) setVolume(volume int) {
	if volume == 0 {
		p.youtube.Call("mute")
	} else {
		p.youtube.Call("unMute")
		p.State.LastNonZeroVolume = volume
	}

	p.youtube.Call("setVolume", volume)
	p.State.Volume = volume
	p.saveState()
	p.Update()

	app.Dispatch(func() {
		app.Window().GetElementByID("volume-var").Set("value", volume)
	})
}

func (p *player) saveState() {
	if err := p.State.save(); err != nil {
		//app.Log("%s", errors.New("saving player state failed").Wrap(err))
		fmt.Println(err)
	}
}

func (p *player) onVolumeChange(ctx app.Context, e app.Event) {
	volume, _ := strconv.Atoi(ctx.JSSrc.Get("value").String())
	p.setVolume(volume)
}

func (p *player) onMute(ctx app.Context, e app.Event) {
	p.setVolume(0)
}

func (p *player) onUnMute(ctx app.Context, e app.Event) {
	p.setVolume(p.State.LastNonZeroVolume)
}

type playerState struct {
	Volume            int
	LastNonZeroVolume int
	Muted             bool
}

func (s *playerState) load() error {
	if err := app.LocalStorage.Get("player.state", s); err != nil {
		//return errors.New("getting player status from local storage failed").
		//	Wrap(err)
		return err
	}

	if *s == (playerState{}) {
		s.Volume = 100
		s.LastNonZeroVolume = 100
	}

	return nil
}

func (s *playerState) save() error {
	if err := app.LocalStorage.Set("player.state", s); err != nil {
		//return errors.New("saving player status in local storage failed").
		//	Wrap(err)
		return err
	}

	return nil
}

func (p *player) Render() app.UI {
	hide := "hide"
	if p.ready {
		hide = ""
	}

	blur := "blur"
	if p.playing {
		blur = ""
	}

	return app.Div().
		ID("player").
		Class("player").
		Body(
			app.Div().
				Class("video").
				Class(blur).
				Body(
					app.Script().
						Src("//www.youtube.com/iframe_api").
						Async(true),
					app.IFrame().
						ID("youtube-player").
						Allow("autoplay").
						Allow("accelerometer").
						Allow("encrypted-media").
						Allow("picture-in-picture").
						Sandbox("allow-presentation allow-same-origin allow-scripts allow-popups").
						Src(fmt.Sprintf(
							"https://www.youtube.com/embed/%s?controls=0&showinfo=0&autoplay=1&loop=1&enablejsapi=1&playsinline=1",
							"uRImXboQnzE",
						)),
				),

			app.Div().
				Class("overlay").
				Body(
					app.Footer().Body(
						app.Stack().
							Class("controls").
							Class(hide).
							Center().
							Content(
								app.If(!p.playing,
									app.Button().
										ID("play").
										Class("button").
										OnClick(p.onPlay).
										Body(
											app.Raw(`
											<svg style="width:24px;height:24px" viewBox="0 0 24 24">
												<path fill="currentColor" d="M8,5.14V19.14L19,12.14L8,5.14Z" />
											</svg>
											`),
										),
								).Else(
									app.Button().
										ID("pause").
										Class("button").
										OnClick(p.onPause).
										Title("Play a random Lofi channel.").
										Body(
											app.Raw(`
											<svg style="width:24px;height:24px" viewBox="0 0 24 24">
    											<path fill="currentColor" d="M14,19.14H18V5.14H14M6,19.14H10V5.14H6V19.14Z" />
											</svg>
											`),
										),
								),
							),
						app.Stack().
							Class("volume").
							Class(hide).
							Center().
							Content(
								app.If(p.State.Volume > 66,
									app.Button().
										Class("button").
										Title("Mute volume").
										OnClick(p.onMute).
										Body(
											app.Raw(`
												<svg style="width:24px;height:24px" viewBox="0 0 24 24">
													<path fill="currentColor" d="M14,3.23V5.29C16.89,6.15 19,8.83 19,12C19,15.17 16.89,17.84 14,18.7V20.77C18,19.86 21,16.28 21,12C21,7.72 18,4.14 14,3.23M16.5,12C16.5,10.23 15.5,8.71 14,7.97V16C15.5,15.29 16.5,13.76 16.5,12M3,9V15H7L12,20V4L7,9H3Z" />
												</svg>
											`),
										),
								).ElseIf(p.State.Volume > 33,
									app.Button().
										Class("button").
										Title("Mute volume").
										OnClick(p.onMute).
										Body(
											app.Raw(`
												<svg style="width:24px;height:24px" viewBox="0 0 24 24">
													<path fill="currentColor" d="M7,9V15H11L16,20V4L11,9H7Z" />
												</svg>
											`),
										),
								).Else(
									app.Button().
										Class("button").
										Title("Unmute volume").
										OnClick(p.onUnMute).
										Body(
											app.Raw(`
												<svg style="width:24px;height:24px" viewBox="0 0 24 24">
													<path fill="currentColor" d="M3,9H7L12,4V20L7,15H3V9M16.59,12L14,9.41L15.41,8L18,10.59L20.59,8L22,9.41L19.41,12L22,14.59L20.59,16L18,13.41L15.41,16L14,14.59L16.59,12Z" />
												</svg>
											`),
										),
								),
								app.Input().
									ID("volume-var").
									Class("volumebar").
									Type("range").
									Placeholder("Volume").
									Min("0").
									Max("100").
									Value(strconv.Itoa(p.State.Volume)).
									OnChange(p.onVolumeChange).
									OnInput(p.onVolumeChange),
							),
					),
				),
		)
}
