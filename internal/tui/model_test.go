package tui

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kwame-Owusu/lista/internal/config"
	"github.com/kwame-Owusu/lista/internal/models"
	"github.com/kwame-Owusu/lista/internal/storage"
	"github.com/muesli/termenv"
)

func TestNewModel(t *testing.T) {
	tl := models.NewTodoList()
	m := NewModel(tl, "test.json")

	if m.todoList != tl {
		t.Error("NewModel() did not store the todo list")
	}

	if m.filename != "test.json" {
		t.Errorf("Expected filename 'test.json', got '%s'", m.filename)
	}

	if m.cursor != 0 {
		t.Errorf("Expected cursor 0, got %d", m.cursor)
	}

	if m.addingTodo {
		t.Error("Expected addingTodo to be false")
	}

	if m.editingTodo {
		t.Error("Expected editingTodo to be false")
	}

	if m.confirmDelete {
		t.Error("Expected confirmDelete to be false")
	}

	if m.confirmPurge {
		t.Error("Expected confirmPurge to be false")
	}

	if m.showHelp {
		t.Error("Expected showHelp to be false")
	}

	if m.focusedField != fieldTitle {
		t.Errorf("Expected focusedField to be fieldTitle, got %v", m.focusedField)
	}

	if m.priority != models.Low {
		t.Errorf("Expected priority %v, got %v", models.Low, m.priority)
	}

	if m.titleInput.Placeholder != "Task title..." {
		t.Errorf("Expected title placeholder 'Task title...', got '%s'", m.titleInput.Placeholder)
	}

	if m.notesInput.Placeholder != "Add notes (optional)..." {
		t.Errorf("Expected notes placeholder 'Add notes (optional)...', got '%s'", m.notesInput.Placeholder)
	}
}

func keyMsg(key string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
}

func TestPurgeCompleted(t *testing.T) {
	tl := models.NewTodoList()
	if err := tl.Add("Done task", models.Low, ""); err != nil {
		t.Fatal(err)
	}
	if err := tl.Add("Pending task", models.Medium, ""); err != nil {
		t.Fatal(err)
	}
	tl.Todos[0].Completed = true

	m := NewModel(tl, filepath.Join(t.TempDir(), "todos.json"))

	// 'c' opens the purge confirmation modal.
	upd, _ := m.Update(keyMsg("c"))
	purgeModel, ok := upd.(model)
	if !ok {
		t.Fatalf("Expected model after Update, got %T", upd)
	}
	if !purgeModel.confirmPurge {
		t.Fatal("Expected confirmPurge to be true after pressing c")
	}
	if !strings.Contains(purgeModel.View(), "Purge all completed todos?") {
		t.Error("Expected purge confirmation text in the rendered view")
	}

	// 'n' cancels without removing anything.
	upd, _ = purgeModel.Update(keyMsg("n"))
	purgeModel, _ = upd.(model)
	if purgeModel.confirmPurge {
		t.Error("Expected confirmPurge to be false after pressing n")
	}
	if purgeModel.todoList.Count() != 2 {
		t.Errorf("Expected 2 todos after cancel, got %d", purgeModel.todoList.Count())
	}

	// 'c' then 'y' purges completed todos.
	upd, _ = purgeModel.Update(keyMsg("c"))
	purgeModel, _ = upd.(model)
	upd, _ = purgeModel.Update(keyMsg("y"))
	purgeModel, _ = upd.(model)
	if purgeModel.confirmPurge {
		t.Error("Expected confirmPurge to be false after confirming")
	}
	if purgeModel.todoList.Count() != 1 {
		t.Errorf("Expected 1 todo after purge, got %d", purgeModel.todoList.Count())
	}
	if purgeModel.todoList.Todos[0].Title != "Pending task" {
		t.Errorf("Expected 'Pending task' to remain, got '%s'", purgeModel.todoList.Todos[0].Title)
	}
}

func TestHelpOverlay(t *testing.T) {
	tl := models.NewTodoList()
	if err := tl.Add("Task", models.Low, ""); err != nil {
		t.Fatal(err)
	}
	m := NewModel(tl, filepath.Join(t.TempDir(), "todos.json"))

	// '?' opens the help overlay.
	upd, _ := m.Update(keyMsg("?"))
	hm, ok := upd.(model)
	if !ok {
		t.Fatalf("Expected model after Update, got %T", upd)
	}
	if !hm.showHelp {
		t.Fatal("Expected showHelp to be true after pressing ?")
	}
	if !strings.Contains(hm.View(), "KEYBINDINGS") {
		t.Error("Expected help overlay to render KEYBINDINGS")
	}

	// Navigation is ignored while the overlay is open.
	upd, _ = hm.Update(keyMsg("k"))
	hm, _ = upd.(model)
	if !hm.showHelp {
		t.Error("Expected help overlay to stay open on non-help keys")
	}
	if hm.cursor != 0 {
		t.Errorf("Expected cursor unchanged while help open, got %d", hm.cursor)
	}

	// 'esc' closes the overlay.
	upd, _ = hm.Update(keyMsg("esc"))
	hm, _ = upd.(model)
	if hm.showHelp {
		t.Error("Expected showHelp to be false after esc")
	}

	// '?' toggles: reopen, then close again.
	upd, _ = hm.Update(keyMsg("?"))
	hm, _ = upd.(model)
	if !hm.showHelp {
		t.Fatal("Expected showHelp to be true after reopening with ?")
	}
	upd, _ = hm.Update(keyMsg("?"))
	hm, _ = upd.(model)
	if hm.showHelp {
		t.Error("Expected showHelp to be false after ? while open")
	}

	// 'q' quits while the overlay is open.
	upd, _ = hm.Update(keyMsg("?"))
	hm, _ = upd.(model)
	if !hm.showHelp {
		t.Fatal("Expected showHelp to be true before quit test")
	}
	if _, cmd := hm.Update(keyMsg("q")); cmd == nil {
		t.Error("Expected a quit command after pressing q while help open")
	}

	// '?' is ignored while a confirmation modal is active.
	upd, _ = hm.Update(keyMsg("?"))
	hm, _ = upd.(model)
	upd, _ = hm.Update(keyMsg("d"))
	hm, _ = upd.(model)
	if !hm.confirmDelete {
		t.Fatal("Expected confirmDelete to be true after pressing d")
	}
	upd, _ = hm.Update(keyMsg("?"))
	hm, _ = upd.(model)
	if hm.showHelp {
		t.Error("Expected help to stay closed while a confirmation modal is active")
	}
}

func TestPriorityMapping(t *testing.T) {
	for _, p := range priorityOptions {
		if got := priorityAt(priorityIndex(p)); got != p {
			t.Errorf("priorityAt(priorityIndex(%v)) = %v, want %v", p, got, p)
		}
	}
}

func TestCyclePriority(t *testing.T) {
	m := NewModel(models.NewTodoList(), "test.json")

	if m.priority != models.Low {
		t.Fatalf("Expected starting priority Low, got %v", m.priority)
	}

	m.cyclePriority(false)
	if m.priority != models.Medium {
		t.Errorf("Expected Medium after cycling down, got %v", m.priority)
	}

	m.cyclePriority(false)
	if m.priority != models.High {
		t.Errorf("Expected High after second cycle down, got %v", m.priority)
	}

	m.cyclePriority(false)
	if m.priority != models.Low {
		t.Errorf("Expected wrap to Low, got %v", m.priority)
	}

	m.cyclePriority(true)
	if m.priority != models.High {
		t.Errorf("Expected High after cycling up from Low, got %v", m.priority)
	}
}

func TestPriorityKeysInForms(t *testing.T) {
	tl := models.NewTodoList()
	if err := tl.Add("Task", models.Low, ""); err != nil {
		t.Fatal(err)
	}
	m := NewModel(tl, filepath.Join(t.TempDir(), "todos.json"))

	// h/l cycle priority in the add form when the priority field is focused.
	upd, _ := m.Update(keyMsg("a"))
	am, _ := upd.(model)
	if !am.addingTodo {
		t.Fatal("Expected addingTodo to be true after pressing a")
	}
	upd, _ = am.Update(keyMsg("tab"))
	am, _ = upd.(model)
	if am.focusedField != fieldPriority {
		t.Fatalf("Expected focus on priority after tab, got %v", am.focusedField)
	}

	upd, _ = am.Update(keyMsg("l"))
	am, _ = upd.(model)
	if am.priority != models.Medium {
		t.Errorf("Expected Medium after l, got %v", am.priority)
	}
	upd, _ = am.Update(keyMsg("h"))
	am, _ = upd.(model)
	if am.priority != models.Low {
		t.Errorf("Expected Low after h, got %v", am.priority)
	}

	// h/l cycle priority in the edit form too.
	upd, _ = m.Update(keyMsg("e"))
	em, _ := upd.(model)
	if !em.editingTodo {
		t.Fatal("Expected editingTodo to be true after pressing e")
	}
	upd, _ = em.Update(keyMsg("tab"))
	em, _ = upd.(model)
	upd, _ = em.Update(keyMsg("l"))
	em, _ = upd.(model)
	if em.priority != models.Medium {
		t.Errorf("Expected Medium after l in edit form, got %v", em.priority)
	}
	upd, _ = em.Update(keyMsg("h"))
	em, _ = upd.(model)
	if em.priority != models.Low {
		t.Errorf("Expected Low after h in edit form, got %v", em.priority)
	}

	// h/l are inserted as text while a text field is focused.
	upd, _ = m.Update(keyMsg("a"))
	am, _ = upd.(model)
	upd, _ = am.Update(keyMsg("h"))
	am, _ = upd.(model)
	upd, _ = am.Update(keyMsg("l"))
	am, _ = upd.(model)
	if got := am.titleInput.Value(); got != "hl" {
		t.Errorf("Expected title input 'hl', got %q", got)
	}
}

func TestTickRefresh(t *testing.T) {
	tl := models.NewTodoList()
	err := tl.Add("Test todo", models.Low, "")
	if err != nil {
		t.Fatal(err)
	}
	tl.Todos[0].CreatedAt = time.Now().Add(-2 * time.Second)

	m := NewModel(tl, "test.json")

	if space := strings.TrimSpace(m.View()); strings.Contains(space, "No todos yet") {
		t.Fatal("Expected rendered list, got 'No todos yet'")
	}

	updated, cmd := m.Update(tickMsg{})
	if cmd == nil {
		t.Error("Expected Update to re-arm the tick command")
	}

	view := updated.View()
	if !strings.Contains(view, "added 2s ago") {
		t.Errorf("Expected fresh timestamp 'added 2s ago' in view, got:\n%s", view)
	}

	updatedForm, cmdForm := updated.Update(tickMsg{})
	if cmdForm == nil {
		t.Error("Expected heartbeat to stay armed while in a form")
	}
	if _, ok := updatedForm.(model); !ok {
		t.Errorf("Expected model back after tick, got %T", updatedForm)
	}
}

func TestTimestampOnlyOnSelectedRow(t *testing.T) {
	tl := models.NewTodoList()
	for _, task := range []string{"First task", "Second task"} {
		if err := tl.Add(task, models.Low, ""); err != nil {
			t.Fatal(err)
		}
	}
	tl.Todos[0].CreatedAt = time.Now().Add(-2 * time.Second)
	tl.Todos[1].CreatedAt = time.Now().Add(-10 * time.Second)

	m := NewModel(tl, "test.json")

	// Cursor starts on the first row, so only its timestamp is visible.
	if view := m.View(); !strings.Contains(view, "added 2s ago") {
		t.Errorf("Expected selected row timestamp 'added 2s ago' in view:\n%s", view)
	} else if strings.Contains(view, "added 10s ago") {
		t.Errorf("Expected non-selected row timestamp to be hidden, got:\n%s", view)
	}

	// Move to the second row: its timestamp appears, the first row's is hidden.
	upd, _ := m.Update(keyMsg("j"))
	moveModel, ok := upd.(model)
	if !ok {
		t.Fatalf("Expected model after Update, got %T", upd)
	}
	if view := moveModel.View(); !strings.Contains(view, "added 10s ago") {
		t.Errorf("Expected selected row timestamp 'added 10s ago' in view:\n%s", view)
	} else if strings.Contains(view, "added 2s ago") {
		t.Errorf("Expected non-selected row timestamp to be hidden, got:\n%s", view)
	}
}

func TestBadgeStrikethroughOnCompleted(t *testing.T) {
	lipgloss.SetColorProfile(termenv.ANSI)
	defer lipgloss.SetColorProfile(termenv.Ascii)

	struck := GetPriorityStyle("High").Strikethrough(true).Render("[High]")
	plain := GetPriorityStyle("High").Strikethrough(false).Render("[High]")

	if !strings.Contains(struck, "\x1b[9m") {
		t.Errorf("Expected completed badge to render with strikethrough, got %q", struck)
	}
	if strings.Contains(plain, "\x1b[9m") {
		t.Errorf("Expected pending badge to render without strikethrough, got %q", plain)
	}
}

func TestCompletedBadgeMuted(t *testing.T) {
	InitStyles(config.DefaultTheme())
	lipgloss.SetColorProfile(termenv.TrueColor)
	defer lipgloss.SetColorProfile(termenv.Ascii)

	muted := getBadgeStyle("High", true).Render("[High]")
	colored := getBadgeStyle("High", false).Render("[High]")

	// Extract the muted foreground SGR token and require it on completed
	// badges while pending badges keep their priority color.
	mutedRef := lipgloss.NewStyle().Foreground(fgMuted).Render("x")
	const sgr = "\x1b[38;2;"
	i := strings.Index(mutedRef, sgr)
	if i < 0 {
		t.Fatalf("Expected muted SGR in rendered style, got %q", mutedRef)
	}
	want := mutedRef[:i+strings.Index(mutedRef[i:], "m")+1]

	if !strings.Contains(muted, want[2:]) {
		t.Errorf("Expected muted foreground %q on completed badge, got %q", want, muted)
	}
	if strings.Contains(colored, want[2:]) {
		t.Errorf("Expected pending badge to keep priority color, got %q", colored)
	}
}

func TestSaveTodosCmd_SnapshotsOnCall(t *testing.T) {
	tl := models.NewTodoList()
	for _, task := range []struct {
		title    string
		priority models.Priority
		notes    string
	}{
		{"Buy groceries", models.Low, ""},
		{"Walk the dog", models.Medium, ""},
		{"Read", models.High, "chapter 3"},
	} {
		if err := tl.Add(task.title, task.priority, task.notes); err != nil {
			t.Fatal(err)
		}
	}

	m := NewModel(tl, filepath.Join(t.TempDir(), "todos.json"))
	cmd := m.saveTodosCmd()

	// Mutate the live list after the save command was created.
	if err := tl.Toggle(1); err != nil {
		t.Fatal(err)
	}
	if err := tl.Add("Should not appear", models.Low, ""); err != nil {
		t.Fatal(err)
	}
	if err := tl.Delete(2); err != nil {
		t.Fatal(err)
	}

	msg := cmd()
	if msg == nil {
		t.Fatal("Expected a msgTodoSaved back from the save command")
	}
	if savedMsg, ok := msg.(msgTodoSaved); !ok {
		t.Fatalf("Expected msgTodoSaved, got %T", msg)
	} else if savedMsg.err != nil {
		t.Fatalf("Save failed: %v", savedMsg.err)
	}

	saved, err := storage.LoadTodos(m.filename)
	if err != nil {
		t.Fatalf("Loading saved file: %v", err)
	}
	if len(saved) != 3 {
		t.Fatalf("Snapshot should have 3 todos, got %d (live list was mutated after save cmd creation)", len(saved))
	}
	if saved[0].Completed {
		t.Error("Snapshot should reflect state at save command creation, not later toggles")
	}
}

func TestSaveTodosCmd_ConcurrentSaves(t *testing.T) {
	const todoCount = 32
	const iterations = 300

	tl := models.NewTodoList()
	for i := 0; i < todoCount; i++ {
		if err := tl.Add(fmt.Sprintf("todo %d", i), models.Low, ""); err != nil {
			t.Fatal(err)
		}
	}

	m := NewModel(tl, filepath.Join(t.TempDir(), "todos.json"))

	jobs := make(chan tea.Cmd)
	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for cmd := range jobs {
				cmd()
			}
		}()
	}

	for i := 0; i < iterations; i++ {
		if err := tl.Toggle(1 + i%todoCount); err != nil {
			t.Fatal(err)
		}
		jobs <- m.saveTodosCmd()
	}
	close(jobs)
	wg.Wait()
}

func executeCmd(t *testing.T, cmd tea.Cmd) {
	t.Helper()
	if cmd == nil {
		t.Fatal("Expected a command to execute")
	}
	if msg := cmd(); msg == nil {
		t.Fatal("Expected a msgTodoSaved back from the save command")
	}
}

func TestUndoToggle(t *testing.T) {
	tl := models.NewTodoList()
	if err := tl.Add("Task", models.Low, ""); err != nil {
		t.Fatal(err)
	}

	file := filepath.Join(t.TempDir(), "todos.json")
	m := NewModel(tl, file)

	// Space completes the todo.
	upd, cmd := m.Update(keyMsg(" "))
	tm, _ := upd.(model)
	if !tm.todoList.Todos[0].Completed {
		t.Fatal("Expected todo to be completed after space")
	}
	executeCmd(t, cmd)

	// u undoes the toggle back to pending.
	upd, cmd = tm.Update(keyMsg("u"))
	tm, _ = upd.(model)
	if tm.todoList.Todos[0].Completed {
		t.Error("Expected todo to be pending after undo")
	}
	executeCmd(t, cmd)

	saved, err := storage.LoadTodos(file)
	if err != nil {
		t.Fatalf("Loading saved file: %v", err)
	}
	if len(saved) != 1 || saved[0].Completed {
		t.Error("Expected an uncompleted todo persisted after undo")
	}
	if len(tm.undoStack) != 0 {
		t.Error("Expected undo stack to be empty after undoing the only action")
	}
}

func TestUndoDelete(t *testing.T) {
	tl := models.NewTodoList()
	if err := tl.Add("To keep", models.Low, ""); err != nil {
		t.Fatal(err)
	}
	if err := tl.Add("To delete", models.Medium, ""); err != nil {
		t.Fatal(err)
	}
	secondID := tl.Todos[1].ID
	nextID := tl.NextID

	file := filepath.Join(t.TempDir(), "todos.json")
	m := NewModel(tl, file)

	// Move cursor to the second todo and delete it.
	m.cursor = 1
	upd, _ := m.Update(keyMsg("d"))
	tm, _ := upd.(model)
	if !tm.confirmDelete {
		t.Fatal("Expected confirmDelete to be true after pressing d")
	}
	upd, cmd := tm.Update(keyMsg("y"))
	tm, _ = upd.(model)
	executeCmd(t, cmd)
	if tm.todoList.Count() != 1 {
		t.Fatalf("Expected 1 todo after delete, got %d", tm.todoList.Count())
	}
	if _, err := tm.todoList.GetByID(secondID); err == nil {
		t.Fatal("Expected deleted todo to be gone")
	}

	// u restores it with the same ID.
	upd, cmd = tm.Update(keyMsg("u"))
	tm, _ = upd.(model)
	executeCmd(t, cmd)
	if tm.todoList.Count() != 2 {
		t.Fatalf("Expected 2 todos after undo, got %d", tm.todoList.Count())
	}
	restored, err := tm.todoList.GetByID(secondID)
	if err != nil {
		t.Fatalf("Expected deleted todo to be restored: %v", err)
	}
	if restored.Title != "To delete" {
		t.Errorf("Expected restored title 'To delete', got '%s'", restored.Title)
	}
	if tm.todoList.NextID != nextID {
		t.Errorf("Expected NextID %d after undo, got %d", nextID, tm.todoList.NextID)
	}
	if tm.cursor != 1 {
		t.Errorf("Expected cursor back at 1 after undo, got %d", tm.cursor)
	}
}

func TestUndoEdit(t *testing.T) {
	tl := models.NewTodoList()
	if err := tl.Add("Original", models.Low, ""); err != nil {
		t.Fatal(err)
	}

	file := filepath.Join(t.TempDir(), "todos.json")
	m := NewModel(tl, file)

	// Open the edit form and change the title.
	upd, _ := m.Update(keyMsg("e"))
	em, _ := upd.(model)
	if !em.editingTodo {
		t.Fatal("Expected editingTodo to be true after pressing e")
	}
	em.titleInput.SetValue("Changed")
	upd, cmd := em.Update(keyMsg("ctrl+s"))
	tm, _ := upd.(model)
	executeCmd(t, cmd)
	if tm.todoList.Todos[0].Title != "Changed" {
		t.Fatalf("Expected title 'Changed' after edit, got '%s'", tm.todoList.Todos[0].Title)
	}

	// u restores the original title.
	upd, cmd = tm.Update(keyMsg("u"))
	tm, _ = upd.(model)
	executeCmd(t, cmd)
	if tm.todoList.Todos[0].Title != "Original" {
		t.Errorf("Expected title 'Original' after undo, got '%s'", tm.todoList.Todos[0].Title)
	}
}

func TestUndoEmptyStackIsNoop(t *testing.T) {
	tl := models.NewTodoList()
	if err := tl.Add("Task", models.Low, ""); err != nil {
		t.Fatal(err)
	}

	m := NewModel(tl, filepath.Join(t.TempDir(), "todos.json"))
	upd, cmd := m.Update(keyMsg("u"))
	tm, _ := upd.(model)

	if cmd != nil {
		t.Error("Expected no command when there is nothing to undo")
	}
	if tm.todoList.Todos[0].Title != "Task" {
		t.Error("Expected no state change when undoing with an empty stack")
	}
	if len(tm.undoStack) != 0 {
		t.Error("Expected undo stack to stay empty")
	}
}

func TestUndoStackIsBounded(t *testing.T) {
	tl := models.NewTodoList()
	if err := tl.Add("Task", models.Low, ""); err != nil {
		t.Fatal(err)
	}

	m := NewModel(tl, filepath.Join(t.TempDir(), "todos.json"))

	// Toggle repeatedly so the oldest snapshots are evicted.
	tm := m
	for i := 0; i < maxUndoEntries+1; i++ {
		upd, _ := tm.Update(keyMsg(" "))
		tm, _ = upd.(model)
	}
	if len(tm.undoStack) != maxUndoEntries {
		t.Fatalf("Expected undo stack capped at %d, got %d", maxUndoEntries, len(tm.undoStack))
	}

	// Undoing maxUndoEntries times reverts to the state after the first toggle,
	// which no longer undoable because it was evicted.
	for i := 0; i < maxUndoEntries; i++ {
		upd, _ := tm.Update(keyMsg("u"))
		tm, _ = upd.(model)
	}
	if len(tm.undoStack) != 0 {
		t.Fatalf("Expected undo stack empty after all undos, got %d", len(tm.undoStack))
	}

	// The todo keeps its completed state from the first toggle.
	if !tm.todoList.Todos[0].Completed {
		t.Error("Expected first (evicted) toggle to remain applied")
	}
}

func TestUndoIgnoredWhileModalsOpen(t *testing.T) {
	tl := models.NewTodoList()
	if err := tl.Add("Task", models.Low, ""); err != nil {
		t.Fatal(err)
	}

	m := NewModel(tl, filepath.Join(t.TempDir(), "todos.json"))

	// u is ignored while the delete confirmation modal is open.
	upd, _ := m.Update(keyMsg("d"))
	dm, _ := upd.(model)
	upd, cmd := dm.Update(keyMsg("u"))
	dm, _ = upd.(model)
	if cmd != nil {
		t.Error("Expected no command when undoing while confirm modal is open")
	}
	if !dm.confirmDelete {
		t.Error("Expected confirm modal to stay open")
	}
	if len(dm.undoStack) != 0 {
		t.Error("Expected no undo snapshot pushed while modal is open")
	}

	// u is ignored while the help overlay is open.
	upd, _ = dm.Update(keyMsg("n"))
	hm, _ := upd.(model)
	upd, _ = hm.Update(keyMsg("?"))
	hm, _ = upd.(model)
	if !hm.showHelp {
		t.Fatal("Expected help overlay to be open")
	}
	upd, cmd = hm.Update(keyMsg("u"))
	hm, _ = upd.(model)
	if cmd != nil {
		t.Error("Expected no undo command while help overlay is open")
	}
	if !hm.showHelp {
		t.Error("Expected help overlay to stay open")
	}
}
