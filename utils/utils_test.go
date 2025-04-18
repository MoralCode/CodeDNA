package utils

import (
	"testing"
)

func TestGetLongestPrefix(t *testing.T) {
	ina := "abcdeefg"
	inb := "abcdegfe"

	want := "abcde"
	msg := GetLongestPrefix(ina, inb)
	if want != msg {
		t.Errorf(`GetLongestPrefix() = %q, want %q`, msg, want)
	}
}

func TestGetLongestPrefixMassive(t *testing.T) {
	ina := "b7e1a36fc7d3322703117c6b5bf6cd285dec59d6202c5fbff2ea1c2f1994dcb42ea8f1de4d3369d4361188bf7b05a336e4a1885f1323a1e9d8e18db2bebf6ad03ae1406b077349b130e3271a05aabe8195163c6f1eeb9c103deb0002984ed19356ae500f759e7d18d1d91f667c58257c331ea3ba985442c7de18f33eb089e77ff45c45e74bca3d6673f6b006d7dc81b571e88b8f03634c8b37bbf0e15a7b924d42c4f45c4aebac17c777628699fc6b5b982c771775776ed09f6e09580376432283da1f6581dfeebe748ecbf22337f93025205d087286d42ed136b49fb225d6e63800300b0ab40a55366428a54a7a8e4a7e34e6893310ebe447a48e8b20e902123b1b5ace6caf7ed376827adfd8007eed2dfb3c8fb8eb10dcd62c4492a4ac310641dfd0f31ce5f60dadc1dd4d3b6b2bbef8abf0509979e46222841f5dca6302ca07b7411db2beb59ef36624da94c039d9c5dc392168edc35aff42dd181073f9a619985ce3c99edee59ddaeffb7c11b6fd08ea7ec27a333d856da2ed90eecc05cb79af6cf92695cd87bc946f86551aa276345c69b3052135cad2b724027ed0ed234bd567c541edb7e789bbe61e1a96aaab0a2f94df4b6c2571171104441f7903f652ad293c530d9c157b88c4507132a194e16fc4328434405e572c89ebb66ff03b66099e7fdc001b049defd2492807f365decde09f50e40c68015c5a74d424d9eba86d5d6cfcd3640c34ce606508f1992b672a88966611be140d1b34270df04b2ea47f609ee9fa7f31932e19641d516fe33a8ff08c2681fd3f1f138eff06103a0349ac66f3e7bef7cb5802f668657bd54095c5d5023d97f27816252788b741bd822346d31a63257a1abd2df672dc6db59f2bac5838279"
	inb := "b7e1a36fc7d3322703117c6b5bf6cd285dec59d6202c5fbff2ea1c2f1994dcb42ea8f1de4d3369d4361188bf7b05a336e4a1885f1323a1e9d8e18db2bebf6ad03ae1406b077349b130e3271a05aabe8195163c6f1eeb9c103deb0002984ed19356ae500f759e7d18d1d91f667c58257c331ea3ba985442c7de18f33eb089e77ff45c45e74bca3d6673f6b006d7dc81b571e88b8f03634c8b37bbf0e15a7b924d42c4f45c4aebac17c777628699fc6b5b982c771775776ed09f6e09580376432283da1f6581dfeebe748ecbf22337f93025205d087286d42ed136b49fb225d6e63800300b0ab40a55366428a54a7a8e4a7e34e6893310ebe447a48e8b20e902123b1b5ace6caf7ed376827adfd8007eed2dfb3c8fb8eb10dcd62c4492a4ac310641dfd0f31ce5f60dadc1dd4d3b6b2bbef8abf0509979e46222841f5dca6302ca07b7411db2beb59ef36624da94c039d9c5dc392168edc35aff42dd181073f9a619985ce3c99edee59ddaeffb7c11b6fd08ea7ec27a333d856da2ed90eecc05cb79af6cf92695cd87bc946f86551aa276345c69b3052135cad2b724027ed0ed234bd567c541edb7e789bbe61e1a96aaab0a2f94df4b6c2571171104441f7903f652ad293c530d9c157b88c4507132a194e16fc4328434405e572c89ebb66ff03b66099e7fdc001b049defd2492807f365decde09f50e40c68015c5a74d424d9eba86d5d6cfcd3640c34ce606508f1992b672a88966611be140d1b34270df04b2ea47f609ee9fa7f31932e19641d516fe33a8ff08c2681fd3f1f138eff06103a0349ac66f3e7bef7cb5802f668657bd54095c5d5023d97f27816252788b741bd822346d31a63257a1abd2df672dc6db59f2bac5838280"

	want := "b7e1a36fc7d3322703117c6b5bf6cd285dec59d6202c5fbff2ea1c2f1994dcb42ea8f1de4d3369d4361188bf7b05a336e4a1885f1323a1e9d8e18db2bebf6ad03ae1406b077349b130e3271a05aabe8195163c6f1eeb9c103deb0002984ed19356ae500f759e7d18d1d91f667c58257c331ea3ba985442c7de18f33eb089e77ff45c45e74bca3d6673f6b006d7dc81b571e88b8f03634c8b37bbf0e15a7b924d42c4f45c4aebac17c777628699fc6b5b982c771775776ed09f6e09580376432283da1f6581dfeebe748ecbf22337f93025205d087286d42ed136b49fb225d6e63800300b0ab40a55366428a54a7a8e4a7e34e6893310ebe447a48e8b20e902123b1b5ace6caf7ed376827adfd8007eed2dfb3c8fb8eb10dcd62c4492a4ac310641dfd0f31ce5f60dadc1dd4d3b6b2bbef8abf0509979e46222841f5dca6302ca07b7411db2beb59ef36624da94c039d9c5dc392168edc35aff42dd181073f9a619985ce3c99edee59ddaeffb7c11b6fd08ea7ec27a333d856da2ed90eecc05cb79af6cf92695cd87bc946f86551aa276345c69b3052135cad2b724027ed0ed234bd567c541edb7e789bbe61e1a96aaab0a2f94df4b6c2571171104441f7903f652ad293c530d9c157b88c4507132a194e16fc4328434405e572c89ebb66ff03b66099e7fdc001b049defd2492807f365decde09f50e40c68015c5a74d424d9eba86d5d6cfcd3640c34ce606508f1992b672a88966611be140d1b34270df04b2ea47f609ee9fa7f31932e19641d516fe33a8ff08c2681fd3f1f138eff06103a0349ac66f3e7bef7cb5802f668657bd54095c5d5023d97f27816252788b741bd822346d31a63257a1abd2df672dc6db59f2bac58382"
	msg := GetLongestPrefix(ina, inb)
	if want != msg {
		t.Errorf(`GetLongestPrefix() = %q, want %q`, msg, want)
	}
}

func TestTrimStringsToEqualLength(t *testing.T) {
	long := "abcdefgh"
	short := "bcdef"
	// same length
	ina, inb := short, short
	a, b, c, d := TrimStringsToEqualLength(ina, inb)
	ex_a, ex_b, ex_c, ex_d := short, short, "", ""

	if a != ex_a || b != ex_b || c != ex_c || d != ex_d {
		t.Errorf(`TestTrimStringsToEqualLength(%q, %q) = %q,%q,%q,%q want %q,%q,%q,%q`, ina, inb, a, b, c, d, ex_a, ex_b, ex_c, ex_d)
	}

	// a longer
	ina, inb = long, short
	a, b, c, d = TrimStringsToEqualLength(ina, inb)
	ex_a, ex_b, ex_c, ex_d = "abcde", short, "fgh", ""

	if a != ex_a || b != ex_b || c != ex_c || d != ex_d {
		t.Errorf(`TestTrimStringsToEqualLength(%q, %q) = %q,%q,%q,%q want %q,%q,%q,%q`, ina, inb, a, b, c, d, ex_a, ex_b, ex_c, ex_d)
	}

	// b longer
	ina, inb = short, long
	a, b, c, d = TrimStringsToEqualLength(ina, inb)
	ex_a, ex_b, ex_c, ex_d = short, "abcde", "", "fgh"

	if a != ex_a || b != ex_b || c != ex_c || d != ex_d {
		t.Errorf(`TestTrimStringsToEqualLength(%q, %q) = %q,%q,%q,%q want %q,%q,%q,%q`, ina, inb, a, b, c, d, ex_a, ex_b, ex_c, ex_d)
	}

}
