import { createClient } from '@connectrpc/connect'
import { useMutation, useQuery } from '@connectrpc/connect-query'
import { useQueryClient } from '@tanstack/react-query'
import { AuthService } from '@todo-app/api-client/src/auth/v1/auth_pb.js'
import { TodoService } from '@todo-app/api-client/src/todo/v1/todo_pb.js'
import { router } from 'expo-router'
import { useState } from 'react'
import { Button, FlatList, Switch, Text, TextInput, TouchableOpacity, View } from 'react-native'
import { storage } from '../../lib/storage'
import { transport } from '../../lib/transport'

export default function TodosScreen() {
  const queryClient = useQueryClient()
  const [newTitle, setNewTitle] = useState('')

  // 一覧取得
  const { data, isLoading } = useQuery(TodoService.method.listTodos, {})

  // 作成
  const createMutation = useMutation(TodoService.method.createTodo, {
    onSuccess: () => {
      queryClient.invalidateQueries()
      setNewTitle('')
    },
  })

  // 更新（completed トグル）
  const updateMutation = useMutation(TodoService.method.updateTodo, {
    onSuccess: () => queryClient.invalidateQueries(),
  })

  // 削除
  const deleteMutation = useMutation(TodoService.method.deleteTodo, {
    onSuccess: () => queryClient.invalidateQueries(),
  })

  // ログアウト
  const handleSignOut = async () => {
    const client = createClient(AuthService, transport)
    await client.signOut({})
    await storage.deleteToken()
    router.replace('/(auth)/login')
  }

  if (isLoading) return <Text>Loading...</Text>

  return (
    <View style={{ flex: 1, padding: 16 }}>
      <View style={{ flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center' }}>
        <Text style={{ fontSize: 24, fontWeight: 'bold' }}>Todos</Text>
        <Button title="Sign Out" onPress={handleSignOut} />
      </View>

      {/* 新規追加 */}
      <View style={{ flexDirection: 'row', marginVertical: 12, gap: 8 }}>
        <TextInput
          style={{ flex: 1, borderWidth: 1, borderColor: '#ccc', padding: 8, borderRadius: 4 }}
          value={newTitle}
          onChangeText={setNewTitle}
          placeholder="New todo..."
          onSubmitEditing={() => {
            if (newTitle.trim()) createMutation.mutate({ title: newTitle })
          }}
        />
        <Button
          title="Add"
          onPress={() => {
            if (newTitle.trim()) createMutation.mutate({ title: newTitle })
          }}
          disabled={createMutation.isPending}
        />
      </View>

      {/* Todo 一覧 */}
      <FlatList
        data={data?.todos ?? []}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => (
          <View
            style={{
              flexDirection: 'row',
              alignItems: 'center',
              paddingVertical: 8,
              borderBottomWidth: 1,
              borderBottomColor: '#eee',
            }}
          >
            <Switch
              value={item.completed}
              onValueChange={(completed) =>
                updateMutation.mutate({ id: item.id, title: item.title, completed })
              }
            />
            <Text
              style={{
                flex: 1,
                marginLeft: 8,
                textDecorationLine: item.completed ? 'line-through' : 'none',
                color: item.completed ? '#999' : '#000',
              }}
            >
              {item.title}
            </Text>
            <TouchableOpacity onPress={() => deleteMutation.mutate({ id: item.id })}>
              <Text style={{ color: 'red' }}>Delete</Text>
            </TouchableOpacity>
          </View>
        )}
      />
    </View>
  )
}
