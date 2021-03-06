/*
 *    Copyright 2018 INS Ecosystem
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package resolver

import (
	"fmt"
	"testing"

	"github.com/insolar/insolar/genesis/mock/storage"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/stretchr/testify/assert"
)

type mockChild struct {
	Reference *object.Reference
	parent    object.Parent
}

func (c *mockChild) GetClassID() string {
	return "mockChild"
}

func (c *mockChild) GetReference() *object.Reference {
	return c.Reference
}

func (c *mockChild) GetParent() object.Parent {
	return c.parent
}

var child = &mockChild{}

type mockParent struct {
	Reference      *object.Reference
	ContextStorage storage.Storage
	parent         object.Parent
}

func (p *mockParent) GetParent() object.Parent {
	return p.parent
}

func (p *mockParent) GetClassID() string {
	return "mockParent"
}

func (p *mockParent) GetReference() *object.Reference {
	return p.Reference
}

func (p *mockParent) GetChildStorage() storage.Storage {
	return nil
}

func (p *mockParent) AddChild(child object.Child) (string, error) {
	return "", nil
}

func (p *mockParent) GetChild(key string) (object.Child, error) {
	return child, nil
}

func (p *mockParent) GetContext() []string {
	return []string{}
}

func (p *mockParent) GetContextStorage() storage.Storage {
	return p.ContextStorage
}

type mockParentWithError struct {
	mockParent
}

func (p *mockParentWithError) GetChild(key string) (object.Child, error) {
	return nil, fmt.Errorf("object with record %s does not exist", key)
}

func TestNewChildResolver(t *testing.T) {
	mockParent := &mockParent{}
	mapStorage := NewChildResolver(mockParent)

	assert.Equal(t, &ChildResolver{
		parent: mockParent,
	}, mapStorage)
}

func TestChildResolver_GetObject_No_Object(t *testing.T) {
	mockParent := &mockParentWithError{}
	resolver := NewChildResolver(mockParent)
	ref, _ := object.NewReference("1", "1", object.ChildScope)

	obj, err := resolver.GetObject(ref, "someClass")

	assert.Equal(t, "object with record 1 does not exist", err.Error())
	assert.Nil(t, obj)
}

func TestChildResolver_GetObject_Wrong_classID(t *testing.T) {
	mockParent := &mockParent{}
	resolver := NewChildResolver(mockParent)
	ref, _ := object.NewReference("1", "1", object.ChildScope)

	obj, err := resolver.GetObject(ref, "someClass")

	assert.Equal(t, "instance class is not `someClass`", err.Error())
	assert.Nil(t, obj)
}

func TestChildResolver_GetObject(t *testing.T) {
	mockParent := &mockParent{}
	resolver := NewChildResolver(mockParent)
	ref, _ := object.NewReference("1", "1", object.ChildScope)

	obj, err := resolver.GetObject(ref, "mockChild")

	assert.Nil(t, err)
	assert.Equal(t, child, obj)
}
