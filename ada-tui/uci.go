package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/WilliamDann/AdaEngine/ada-chess/core"
	"github.com/WilliamDann/AdaEngine/ada-chess/movegen"
	"github.com/WilliamDann/AdaEngine/ada-chess/position"
)

// uciEngine manages a UCI engine subprocess.
type uciEngine struct {
	cmd    *exec.Cmd
	stdin  *bufio.Writer
	stdout *bufio.Scanner
}

func startUCI(path string) (*uciEngine, error) {
	cmd := exec.Command(path)
	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	e := &uciEngine{
		cmd:    cmd,
		stdin:  bufio.NewWriter(in),
		stdout: bufio.NewScanner(out),
	}

	e.send("uci")
	e.waitFor("uciok")
	e.send("setoption name Threads value 8")
	e.send("isready")
	e.waitFor("readyok")

	return e, nil
}

func (e *uciEngine) send(line string) {
	fmt.Fprintln(e.stdin, line)
	e.stdin.Flush()
}

func (e *uciEngine) waitFor(prefix string) string {
	for e.stdout.Scan() {
		line := e.stdout.Text()
		if strings.HasPrefix(line, prefix) {
			return line
		}
	}
	return ""
}

// bestMove asks the engine for a move. moveHistory is a list of UCI move
// strings (e.g. "e2e4") from the starting position.
func (e *uciEngine) bestMove(startFEN string, moveHistory []string, thinkTime time.Duration) string {
	posCmd := "position fen " + startFEN
	if len(moveHistory) > 0 {
		posCmd += " moves " + strings.Join(moveHistory, " ")
	}
	e.send(posCmd)
	e.send(fmt.Sprintf("go movetime %d", thinkTime.Milliseconds()))

	for e.stdout.Scan() {
		line := e.stdout.Text()
		if strings.HasPrefix(line, "bestmove") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1]
			}
		}
	}
	return ""
}

func (e *uciEngine) close() {
	e.send("quit")
	e.cmd.Wait()
}

// parseSFMove converts a UCI move string (e.g. "e2e4", "e7e8q") to a
// core.Move by matching it against the legal moves in the position.
func parseSFMove(pos *position.Position, uci string) (core.Move, bool) {
	moves := movegen.LegalMoves(pos)
	for i := 0; i < moves.Count(); i++ {
		m := moves.Get(i)
		if m.String() == uci {
			return m, true
		}
	}
	return core.NoMove, false
}
