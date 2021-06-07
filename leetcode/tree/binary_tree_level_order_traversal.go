package tree

/**
 * Definition for a binary tree node.
 */
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func levelOrder(root *TreeNode) [][]int {

	result := [][]int{}
	iterator(root, &result, 1)
	return result
}

func iterator(root *TreeNode, result *[][]int, level int) {

	if root == nil {
		return
	}

	if len(*result) < level {
		rootVal := []int{root.Val}
		*result = append(*result, rootVal)
	} else {
		(*result)[level-1] = append((*result)[level-1], root.Val)
	}

	iterator(root.Left, result, level+1)
	iterator(root.Right, result, level+1)

	return
}
