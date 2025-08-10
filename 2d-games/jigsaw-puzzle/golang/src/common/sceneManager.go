package common

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Update(context *SceneContext) error
	Draw(screen *ebiten.Image, context *SceneContext)
}

type SceneContext struct {
	SceneManager *SceneManager
}

type SceneManager struct {
	current      Scene
	sceneContext *SceneContext
	scenes       map[string]func() Scene
}

func NewSceneManager() *SceneManager {
	m := &SceneManager{
		sceneContext: &SceneContext{},
		scenes:       map[string]func() Scene{},
	}

	m.sceneContext.SceneManager = m
	return m
}

func (s *SceneManager) AddScene(name string, newSceneFunc func() Scene) {
	s.scenes[name] = newSceneFunc
}

func (s *SceneManager) SetScene(name string) {
	newSceneFunc, b := s.scenes[name]
	if !b {
		panic(fmt.Sprintf("%s scene not found", name))
	}

	s.current = newSceneFunc()
}

func (s *SceneManager) Update() error {
	if s.current != nil {
		return s.current.Update(s.sceneContext)
	}

	return nil
}

func (s *SceneManager) Draw(screen *ebiten.Image) {
	if s.current != nil {
		s.current.Draw(screen, s.sceneContext)
	}
}
