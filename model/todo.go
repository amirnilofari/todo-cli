package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexeyco/simpletable"
	"io/ioutil"
	"os"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Todes []item

func (t *Todes) Add(task string) {

	todo := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*t = append(*t, todo)
}

func (t *Todes) Complete(index int) error {
	ls := *t

	if index <= 0 || index > len(ls) {
		return errors.New("invalid index")
	}

	ls[index-1].CompletedAt = time.Now()
	ls[index-1].Done = true

	return nil
}

func (t *Todes) Delete(index int) error {
	ls := *t

	if index <= 0 || index > len(ls) {
		return errors.New("invalid index")
	}

	*t = append(ls[:index-1], ls[index:]...)

	return nil
}

func (t *Todes) Load(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return err
	}

	err = json.Unmarshal(file, t)
	if err != nil {
		return err
	}

	return nil

}

func (t *Todes) Store(filename string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func (t *Todes) Print() {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done?"},
			{Align: simpletable.AlignCenter, Text: "CreatedAt"},
			{Align: simpletable.AlignCenter, Text: "CompletedAt"},
		},
	}

	var cells [][]*simpletable.Cell

	for idx, item := range *t {
		idx++

		index := Gray(fmt.Sprintf("%d", idx))
		task := Blue(item.Task)
		done := Red("⛔️")
		createdAt := Gray(item.CreatedAt.Format(time.RFC1123))
		completedAt := Gray(item.CompletedAt.Format(time.RFC1123))

		if item.Done {
			task = Green(fmt.Sprintf("\u2705 %s", item.Task))
			done = GreenMark("✅")
		}
		cells = append(cells, *&[]*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: index},
			{Align: simpletable.AlignCenter, Text: task},
			{Align: simpletable.AlignCenter, Text: done},
			{Align: simpletable.AlignCenter, Text: createdAt},
			{Align: simpletable.AlignCenter, Text: completedAt},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 5, Text: fmt.Sprintf("You have %d pending todos", t.CountPending())},
		},
	}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}

func (t *Todes) CountPending() int {
	total := 0

	for _, item := range *t {
		if !item.Done {
			total++
		}
	}
	return total
}
