export async function fetchADOComments(taskId: number) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/commentservice')
  return m.FetchADOComments(taskId)
}

export async function replyToADOComment(taskId: number, content: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/commentservice')
  return m.ReplyToADOComment(taskId, content)
}

export async function listComments(taskId: number) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/commentservice')
  return m.ListComments(taskId)
}

export async function addComment(taskId: number, content: string) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/commentservice')
  return m.AddComment(taskId, content)
}

export async function pushCommentToADO(commentId: number) {
  const m = await import('../../bindings/dev.azure.com/xbox/xb-tasks/internal/app/commentservice')
  return m.PushCommentToADO(commentId)
}
