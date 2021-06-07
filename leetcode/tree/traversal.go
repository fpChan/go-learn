package tree

func inOrder(root *TreeNode, array []int) {
	if root == nil {
		return
	}
	inOrder(root.Left, array)
	array = append(array, root.Val)
	inOrder(root.Right, array)
}
