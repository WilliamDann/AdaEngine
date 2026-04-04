package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/fen"
	"github.com/WilliamDann/AdaEngine/ada-chess/movegen"
	"github.com/WilliamDann/AdaEngine/ada-chess/pgn"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
	"github.com/WilliamDann/AdaEngine/ada-search"
)

const defaultDepth = 4

const (
	modeHuman = 0 // human vs engine (engine replies after human moves)
	modeAuto  = 1 // engine vs engine
	modeOff   = 2 // human vs human (no auto-play)
)

type app struct {
	pos       *position.Position
	depth     int
	threads   int
	timeLimit time.Duration
	mode      int
	game      *pgn.Game

	tv    *tview.Application
	board *KittyImage
	log   *tview.TextView
	input *tview.InputField
	info  *tview.TextView
}

func newApp() *app {
	return &app{
		depth:     defaultDepth,
		threads:   0,
		timeLimit: 20 * time.Second,
	}
}

func (a *app) searchLabel() string {
	if a.timeLimit > 0 {
		return fmt.Sprintf("depth 1..%d, %s limit", a.depth, a.timeLimit)
	}
	return fmt.Sprintf("depth 1..%d", a.depth)
}

func (a *app) appendLog(msg string) {
	fmt.Fprintf(a.log, "%s\n", msg)
	a.log.ScrollToEnd()
}

func (a *app) updateBoard() {
	a.board.SetImage(renderBoard(a.pos))
}

func (a *app) updateInfo() {
	a.info.Clear()
	color := "[aqua]White[-]"
	if a.pos.ActiveColor == core.Black {
		color = "[teal]Black[-]"
	}

	status := ""
	moves := movegen.LegalMoves(a.pos)
	if movegen.InCheck(a.pos) {
		if moves.Count() == 0 {
			status = " [red]CHECKMATE[-]"
		} else {
			status = " [red]CHECK[-]"
		}
	} else if moves.Count() == 0 {
		status = " [yellow]STALEMATE[-]"
	}

	fmt.Fprintf(a.info, " %s to move%s  |  Depth: [aqua]%d[-]  |  Moves: [aqua]%d[-]  |  Move [aqua]%d[-]",
		color, status, a.depth, moves.Count(), a.pos.Fullmoves)
}

func (a *app) refresh() {
	a.updateBoard()
	a.updateInfo()
}

// engineMove starts an engine search and plays the result. If the game isn't
// over and auto mode is on, it schedules another move.
func (a *app) engineMove() {
	moves := movegen.LegalMoves(a.pos)
	if moves.Count() == 0 {
		return
	}
	d := a.depth
	pos := a.pos
	a.appendLog(fmt.Sprintf("[yellow]Thinking (%s)...[-]", a.searchLabel()))
	go func() {
		start := time.Now()
		res := search.Search(pos, d, a.threads, a.timeLimit, func(r search.Result) {
			elapsed := time.Since(start)
			a.tv.QueueUpdateDraw(func() {
				a.appendLog(fmt.Sprintf("  depth [aqua]%d[-]: [aqua]%s[-]  score: [yellow]%s[-]  nodes: %d  time: [yellow]%s[-]",
					r.Depth, r.Move, formatScore(r.Score), r.Nodes, elapsed.Round(time.Millisecond)))
			})
		})
		elapsed := time.Since(start)
		a.tv.QueueUpdateDraw(func() {
			if res.Move == core.NoMove {
				return
			}
			a.game.AddMove(pos, res.Move)
			a.pos = position.MakeMove(pos, res.Move)
			nps := uint64(0)
			if elapsed.Seconds() > 0 {
				nps = uint64(float64(res.Nodes) / elapsed.Seconds())
			}
			a.appendLog(fmt.Sprintf("Engine: [aqua]%s[-]  score: [yellow]%s[-]  nodes: %d  time: [yellow]%s[-]  nps: [yellow]%d[-]",
				res.Move, formatScore(res.Score), res.Nodes, elapsed.Round(time.Millisecond), nps))
			a.refresh()

			// Continue if auto mode or game not over
			if a.mode == modeAuto {
				next := movegen.LegalMoves(a.pos)
				if next.Count() > 0 {
					a.engineMove()
				}
			}
		})
	}()
}

func (a *app) handleInput(text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}

	args := strings.Fields(text)
	switch args[0] {
	case "quit", "exit", "q":
		a.tv.Stop()

	case "help", "h":
		a.appendLog("[aqua]Commands:[-]")
		a.appendLog("  [yellow]<move>[-]       e.g. e2e4, e7e8q")
		a.appendLog("  [yellow]play [d][-]     Engine makes a move")
		a.appendLog("  [yellow]search [d][-]   Show best move")
		a.appendLog("  [yellow]auto[-]         Engine vs engine")
		a.appendLog("  [yellow]mode[-]         Cycle: human/auto/off")
		a.appendLog("  [yellow]stop[-]         Stop auto-play")
		a.appendLog("  [yellow]moves[-]        List legal moves")
		a.appendLog("  [yellow]depth <n>[-]    Set search depth")
		a.appendLog("  [yellow]time <dur>[-]   Set time limit (e.g. 5s, 20s, 0 to disable)")
		a.appendLog("  [yellow]threads <n>[-]  Set search threads")
		a.appendLog("  [yellow]fen <str>[-]    Load position")
		a.appendLog("  [yellow]new[-]          New game")
		a.appendLog("  [yellow]pgn[-]          Show PGN of current game")
		a.appendLog("  [yellow]quit[-]         Exit")

	case "moves", "m":
		moves := movegen.LegalMoves(a.pos)
		var parts []string
		for i := 0; i < moves.Count(); i++ {
			parts = append(parts, moves.Get(i).String())
		}
		a.appendLog(fmt.Sprintf("[yellow]%d moves:[-] %s", len(parts), strings.Join(parts, " ")))

	case "depth", "d":
		if len(args) < 2 {
			a.appendLog(fmt.Sprintf("Depth: [aqua]%d[-]", a.depth))
		} else if d, err := strconv.Atoi(args[1]); err == nil && d > 0 {
			a.depth = d
			a.appendLog(fmt.Sprintf("Depth set to [aqua]%d[-]", d))
			a.updateInfo()
		}

	case "time":
		if len(args) < 2 {
			if a.timeLimit > 0 {
				a.appendLog(fmt.Sprintf("Time limit: [aqua]%s[-]", a.timeLimit))
			} else {
				a.appendLog("Time limit: [aqua]off[-] (depth only)")
			}
		} else if args[1] == "0" || args[1] == "off" {
			a.timeLimit = 0
			a.appendLog("Time limit [aqua]disabled[-] (depth only)")
		} else if dur, err := time.ParseDuration(args[1]); err == nil && dur > 0 {
			a.timeLimit = dur
			a.appendLog(fmt.Sprintf("Time limit set to [aqua]%s[-]", dur))
		} else {
			a.appendLog("[red]Usage: time <duration> (e.g. 5s, 20s, 1m, 0)[-]")
		}

	case "threads", "t":
		if len(args) < 2 {
			t := a.threads
			if t <= 0 {
				t = runtime.NumCPU()
			}
			a.appendLog(fmt.Sprintf("Threads: [aqua]%d[-]", t))
		} else if t, err := strconv.Atoi(args[1]); err == nil && t > 0 {
			a.threads = t
			a.appendLog(fmt.Sprintf("Threads set to [aqua]%d[-]", t))
		}

	case "search", "s":
		d := a.parseDepthArg(args)
		a.appendLog(fmt.Sprintf("[yellow]Searching (%s)...[-]", a.searchLabel()))
		pos := a.pos
		go func() {
			start := time.Now()
			res := search.Search(pos, d, a.threads, a.timeLimit, func(r search.Result) {
				elapsed := time.Since(start)
				a.tv.QueueUpdateDraw(func() {
					a.appendLog(fmt.Sprintf("  depth [aqua]%d[-]: [aqua]%s[-]  score: [yellow]%s[-]  nodes: %d  time: [yellow]%s[-]",
						r.Depth, r.Move, formatScore(r.Score), r.Nodes, elapsed.Round(time.Millisecond)))
				})
			})
			elapsed := time.Since(start)
			a.tv.QueueUpdateDraw(func() {
				if res.Move == core.NoMove {
					a.appendLog("[red]No moves available.[-]")
				} else {
					nps := uint64(0)
					if elapsed.Seconds() > 0 {
						nps = uint64(float64(res.Nodes) / elapsed.Seconds())
					}
					a.appendLog(fmt.Sprintf("Best: [aqua]%s[-]  score: [yellow]%s[-]  nodes: %d  time: [yellow]%s[-]  nps: [yellow]%d[-]",
						res.Move, formatScore(res.Score), res.Nodes, elapsed.Round(time.Millisecond), nps))
				}
			})
		}()

	case "play", "p":
		d := a.parseDepthArg(args)
		a.appendLog(fmt.Sprintf("[yellow]Thinking (%s)...[-]", a.searchLabel()))
		pos := a.pos
		go func() {
			start := time.Now()
			res := search.Search(pos, d, a.threads, a.timeLimit, func(r search.Result) {
				elapsed := time.Since(start)
				a.tv.QueueUpdateDraw(func() {
					a.appendLog(fmt.Sprintf("  depth [aqua]%d[-]: [aqua]%s[-]  score: [yellow]%s[-]  nodes: %d  time: [yellow]%s[-]",
						r.Depth, r.Move, formatScore(r.Score), r.Nodes, elapsed.Round(time.Millisecond)))
				})
			})
			elapsed := time.Since(start)
			a.tv.QueueUpdateDraw(func() {
				if res.Move == core.NoMove {
					a.appendLog("[red]No moves available.[-]")
				} else {
					a.game.AddMove(pos, res.Move)
					a.pos = position.MakeMove(pos, res.Move)
					nps := uint64(0)
					if elapsed.Seconds() > 0 {
						nps = uint64(float64(res.Nodes) / elapsed.Seconds())
					}
					a.appendLog(fmt.Sprintf("Engine: [aqua]%s[-]  score: [yellow]%s[-]  nodes: %d  time: [yellow]%s[-]  nps: [yellow]%d[-]",
						res.Move, formatScore(res.Score), res.Nodes, elapsed.Round(time.Millisecond), nps))
					a.refresh()
				}
			})
		}()

	case "auto":
		a.mode = modeAuto
		a.appendLog("[yellow]Auto-play started.[-] Type [yellow]stop[-] to pause.")
		a.engineMove()

	case "mode":
		a.mode = (a.mode + 1) % 3
		names := []string{"human (engine replies)", "auto (engine vs engine)", "off (manual only)"}
		a.appendLog(fmt.Sprintf("Mode: [aqua]%s[-]", names[a.mode]))

	case "stop":
		a.mode = modeOff
		a.appendLog("[yellow]Auto-play stopped.[-]")

	case "new":
		a.mode = modeHuman
		a.pos, _ = fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
		a.game = pgn.NewGame(a.pos)
		a.log.Clear()
		fmt.Fprint(a.log, logoString())
		a.appendLog("[yellow]New game.[-]\n")
		a.refresh()

	case "fen":
		if len(args) < 2 {
			a.appendLog("[red]Usage: fen <fen-string>[-]")
		} else {
			fenStr := strings.Join(args[1:], " ")
			if p, err := fen.Parse(fenStr); err == nil {
				a.pos = p
				a.game = pgn.NewGame(a.pos)
				a.appendLog("[yellow]Position loaded.[-]")
				a.refresh()
			} else {
				a.appendLog(fmt.Sprintf("[red]Invalid FEN: %v[-]", err))
			}
		}

	case "pgn":
		a.appendLog("[aqua]--- PGN ---[-]")
		a.appendLog(a.game.String())

	default:
		m, ok := parseMove(a.pos, text)
		if ok {
			a.game.AddMove(a.pos, m)
			a.pos = position.MakeMove(a.pos, m)
			a.appendLog(fmt.Sprintf("You: [aqua]%s[-]", m))
			a.refresh()
			// Auto-reply in human mode
			if a.mode == modeHuman {
				next := movegen.LegalMoves(a.pos)
				if next.Count() > 0 {
					a.engineMove()
				}
			}
		} else {
			a.appendLog(fmt.Sprintf("[red]Unknown command or illegal move: %s[-]", text))
		}
	}
}

func (a *app) parseDepthArg(args []string) int {
	d := a.depth
	if len(args) >= 2 {
		if parsed, err := strconv.Atoi(args[1]); err == nil && parsed > 0 {
			d = parsed
		}
	}
	return d
}

func logoString() string {
	lines := []string{
		`     **     *******       **`,
		`    ****   /**////**     ****`,
		`   **//**  /**    /**   **//**`,
		`  **  //** /**    /**  **  //**`,
		` **********/**    /** **********`,
		`/**//////**/**    ** /**//////**`,
		`/**     /**/*******  /**     /**`,
		`//      // ///////   //      //`,
	}
	var sb strings.Builder
	for i, l := range lines {
		for _, ch := range l {
			switch ch {
			case '*':
				sb.WriteString("[aqua]" + string(ch) + "[-]")
			case '/', '\\':
				sb.WriteString("[blue]" + string(ch) + "[-]")
			default:
				sb.WriteRune(ch)
			}
		}
		if i == 3 {
			sb.WriteString("   [white]Chess Engine[-]")
		} else if i == 4 {
			sb.WriteString("  [white]v0.1[-]")
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func formatScore(score int) string {
	if score >= search.Mate-500 {
		return "mate in " + strconv.Itoa(score - search.Mate)
	}
	if score <= -search.Mate+500 {
		return "mate in -" + strconv.Itoa(score + (-search.Mate))
	}
	return fmt.Sprintf("%.2f", float64(score)/100.0)
}

func parseMove(pos *position.Position, input string) (core.Move, bool) {
	input = strings.TrimSpace(input)
	moves := movegen.LegalMoves(pos)
	for i := 0; i < moves.Count(); i++ {
		m := moves.Get(i)
		// Match UCI (e.g. e2e4) or SAN (e.g. e4, Nf3, O-O)
		if m.String() == strings.ToLower(input) || pgn.SAN(pos, m) == input {
			return m, true
		}
	}
	return core.NoMove, false
}

func main() {
	fmt.Print("\033[2J\033[H")
	fmt.Println("Initializing engine...")
	startPos, _ := fen.Parse("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	movegen.LegalMoves(startPos)
	loadPieces()

	tv := tview.NewApplication()
	a := newApp()
	a.pos = startPos
	a.game = pgn.NewGame(startPos)
	a.tv = tv

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault

	// Board (right) — Kitty graphics widget
	a.board = NewKittyImage()
	a.board.SetBackgroundColor(tcell.ColorDefault)

	// Logo + log (left)
	a.log = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	a.log.SetBorder(false).SetBackgroundColor(tcell.ColorDefault)
	fmt.Fprint(a.log, logoString())
	a.appendLog("[aqua]Ready.[-] Type [yellow]help[-] for commands.\n")

	// Info bar
	a.info = tview.NewTextView().
		SetDynamicColors(true)
	a.info.SetBorder(false).SetBackgroundColor(tcell.ColorDefault)

	// Input
	a.input = tview.NewInputField().
		SetLabel("[aqua]> [-]").
		SetLabelColor(tcell.ColorAqua).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetFieldTextColor(tcell.ColorWhite).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				text := a.input.GetText()
				a.input.SetText("")
				a.handleInput(text)
			}
		})
	a.input.SetBackgroundColor(tcell.ColorDefault)

	a.refresh()

	// Layout:
	//  ┌──────────┬──────────┐
	//  │   log    │  board   │
	//  ├──────────┴──────────┤
	//  │      info bar       │
	//  ├─────────────────────┤
	//  │      input          │
	//  └─────────────────────┘
	upper := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(a.log, 0, 1, false).
		AddItem(a.board, 0, 1, false)
	upper.SetBackgroundColor(tcell.ColorDefault)

	root := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(upper, 0, 1, false).
		AddItem(a.info, 1, 0, false).
		AddItem(a.input, 1, 0, true)
	root.SetBackgroundColor(tcell.ColorDefault)

	tv.SetRoot(root, true).SetFocus(a.input)

	if err := tv.Run(); err != nil {
		panic(err)
	}
}
