# MediaSum research-use and citation notice

ISC-RAG uses **MediaSum: A Large-scale Media Interview Dataset for Dialogue
Summarization** by Chenguang Zhu, Yang Liu, Jie Mei, and Michael Zeng.

- Official project: <https://github.com/zcgzcgzcg1/MediaSum>
- Paper: <https://arxiv.org/abs/2103.06410>
- Dataset mirror used by this project: <https://huggingface.co/datasets/nbroad/mediasum>
- Pinned mirror revision: `bcbcef7de3b4e702e526352d825a4ff06de2becb`
- Pinned official repository commit: `7c92589bb6e69e678d043a9cf8c8c7285cd88055`

## Terms carried forward

The MediaSum authors ask users to restrict the dataset to research purposes and
to cite their paper. The official repository does not include a standalone
license file. The Hugging Face mirror currently displays `cc-by-nc-sa-4.0` in
its page metadata, while the dataset card itself says that no license was
specified and repeats the authors' research-only request. This repository
therefore follows the more restrictive, explicit instruction:

- use MediaSum and MediaSum-derived ISC-RAG artifacts only for research;
- cite the MediaSum paper;
- do not treat ISC-RAG's MIT license as a license to the underlying transcripts;
- consult the upstream sources before redistributing corpus content.

The full 4.45 GB corpus is not stored in this Git repository. `./bootstrap.sh`
downloads the pinned upstream files and verifies them against
`resources/mediasum/MANIFEST.sha256`. The pinned official repository is also
available as `sources/MediaSum/`.

Source-derived questions, calibration sheets, illustrative facts, and quoted
evidence retained in this repository exist to document and reproduce the
research. They remain subject to the notice above and to any rights held by the
original NPR/CNN sources.

## Citation

```bibtex
@article{zhu2021mediasum,
  title={MediaSum: A Large-scale Media Interview Dataset for Dialogue Summarization},
  author={Zhu, Chenguang and Liu, Yang and Mei, Jie and Zeng, Michael},
  journal={arXiv preprint arXiv:2103.06410},
  year={2021}
}
```
