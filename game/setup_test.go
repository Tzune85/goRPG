package game

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestSetupEOFName(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	reader := strings.NewReader("1\n")
	game.in = reader
	game.scanner = bufio.NewScanner(reader)

	game.setup()

	if game.running {
		t.Errorf("expected running=false after EOF")
	}
	if strings.Contains(buf.String(), "Farewell") {
		t.Errorf("EOF path should not print quit message")
	}
}

func TestSetupEOFClass(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	reader := strings.NewReader("1\nHero\n")
	game.in = reader
	game.scanner = bufio.NewScanner(reader)

	game.setup()

	if game.running {
		t.Errorf("expected running=false after EOF")
	}
	if strings.Contains(buf.String(), "Farewell") {
		t.Errorf("EOF path should not print quit message")
	}
}

func TestSetupDefaultName(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	reader := strings.NewReader("1\n\n1\n")
	game.in = reader
	game.scanner = bufio.NewScanner(reader)

	game.setup()

	if game.Player.Name != "Hero" {
		t.Errorf("expected hero named Hero")
	}
}

func TestSetupDescribeRoom(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	reader := strings.NewReader("1\n\n1\n")
	game.in = reader
	game.scanner = bufio.NewScanner(reader)

	game.setup()
	output := buf.String()

	if !strings.Contains(output, "threshold of darkness") {
		t.Errorf("expected description, got : %s", output)
	}
}

func TestSetupCreatePlayer(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	reader := strings.NewReader("1\nHero\n1\n")
	game.in = reader
	game.scanner = bufio.NewScanner(reader)

	game.setup()

	if game.Player.Name != "Hero" {
		t.Errorf("expected player 'Hero' created")
	}

	if game.Player.Class != Warrior {
		t.Errorf("expected player 'Warrior' created")
	}

}

func TestChooseClass(t *testing.T) {

	cases := []struct {
		name     string
		input    string
		expected Class
	}{
		{"digit warrior", "1\n", Warrior},
		{"word warrior", "warrior\n", Warrior},
		{"digit mage", "2\n", Mage},
		{"word mage", "mage\n", Mage},
		{"digit thief", "3\n", Thief},
		{"word thief", "thief\n", Thief},
		{"word god", "god\n", God},
		{"digit god", "42\n", God},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			game := New()
			game.out = &buf
			reader := strings.NewReader(tc.input)
			game.in = reader
			game.scanner = bufio.NewScanner(reader)
			class, ok := game.chooseClass()
			if !ok || class != tc.expected {
				t.Errorf("input %q: expected %v, got %v (ok=%v)", tc.input, tc.expected, class, ok)
			}
		})
	}

}

func TestChoseClassInvalid(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	reader := strings.NewReader("Hero\n4\n")
	game.in = reader
	game.scanner = bufio.NewScanner(reader)

	game.chooseClass()
	output := buf.String()

	if !strings.Contains(output, "Please choose 1, 2, or 3.") {
		t.Errorf("expected Please choose, got : %s", output)
	}
}

func TestItTranslator(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	reader := strings.NewReader("2\n")
	game.in = reader
	game.scanner = bufio.NewScanner(reader)

	game.setup()

	output := buf.String()

	if !strings.Contains(output, "Inserisci il tuo nome, avventuriero:") {
		t.Errorf("expected italian language, got : %s", output)
	}

}

func TestEnTranslator(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	reader := strings.NewReader("1\n")
	game.in = reader
	game.scanner = bufio.NewScanner(reader)

	game.setup()

	output := buf.String()

	if !strings.Contains(output, "Enter your name, adventurer: ") {
		t.Errorf("expected english language, got : %s", output)
	}

}

func TestTranslatorEOF(t *testing.T) {
	var buf bytes.Buffer
	game := New()
	game.out = &buf
	reader := strings.NewReader("\n")
	game.in = reader
	game.scanner = bufio.NewScanner(reader)

	game.setup()

	if game.running {
		t.Errorf("expected running=false after EOF")
	}
	if strings.Contains(buf.String(), "Farewell") {
		t.Errorf("EOF path should not print quit message")
	}
}
