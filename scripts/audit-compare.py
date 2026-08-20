#!/usr/bin/env python3
"""Compare freshly regenerated McNemar tables against the committed ones.

Usage: audit_compare.py COMMITTED=REGENERATED [COMMITTED=REGENERATED ...]

Rows are compared positionally, section by section, because the same
(budget, comparison) key legitimately appears in more than one section with
different numbers. Any difference at all is reported: this is an audit, not a
tolerance check.
"""
import sys, re

# Section headers appear at two levels in the committed files: '# McNemar
# paired analysis - ...' for the main families, and '## <free text>' for tables
# added later on request. Both start a new section.
HDR = re.compile(r'^#\s*McNemar paired analysis\s*[\u2014-]\s*(.*)$')
SUBHDR = re.compile(r'^##\s+(?!#)(.*\S)\s*$')

def parse(path):
    """-> list of (section_title, [row_cells, ...])"""
    sections, cur = [], None
    for line in open(path, encoding='utf-8'):
        line = line.rstrip('\n')
        m = HDR.match(line.strip())
        if m:
            cur = (m.group(1).strip(), [])
            sections.append(cur)
            continue
        m = SUBHDR.match(line.strip())
        if m:
            cur = (m.group(1).strip(), [])
            sections.append(cur)
            continue
        if cur is None:
            continue
        s = line.strip()
        if not s.startswith('|'):
            continue
        cells = [c.strip() for c in s.strip('|').split('|')]
        if not cells or cells[0] in ('Budget',) or set(cells[0]) <= set('-: '):
            continue
        cur[1].append(cells)
    return sections

def main():
    total_rows = total_diffs = 0
    problems = []
    for arg in sys.argv[1:]:
        committed_path, regen_path = arg.split('=', 1)
        print('=' * 78)
        print('COMMITTED :', committed_path)
        print('REGENERATED:', regen_path)
        try:
            A, B = parse(committed_path), parse(regen_path)
        except OSError as e:
            problems.append(f'{committed_path}: {e}')
            print('  ERROR:', e); continue

        if len(A) != len(B):
            problems.append(f'{committed_path}: section count {len(A)} vs {len(B)}')
            print(f'  MISMATCH: {len(A)} sections committed, {len(B)} regenerated')

        for i, (sa, sb) in enumerate(zip(A, B), 1):
            ta, ra = sa; tb, rb = sb
            print(f'\n  -- section {i}: {ta}')
            if ta != tb:
                problems.append(f'{committed_path} s{i}: header differs')
                print(f'     HEADER DIFFERS: regenerated = {tb}')
            if len(ra) != len(rb):
                problems.append(f'{committed_path} s{i}: row count {len(ra)} vs {len(rb)}')
                print(f'     ROW COUNT: {len(ra)} committed, {len(rb)} regenerated')
            n = diffs = 0
            for j, (x, y) in enumerate(zip(ra, rb), 1):
                n += 1; total_rows += 1
                if x != y:
                    diffs += 1; total_diffs += 1
                    label = ' / '.join(x[:2]) if len(x) >= 2 else str(j)
                    problems.append(f'{committed_path} s{i} row {j} ({label})')
                    print(f'     ROW {j} DIFFERS  [{label}]')
                    for k, (cx, cy) in enumerate(zip(x, y)):
                        if cx != cy:
                            print(f'        col {k}: committed={cx!r}  regenerated={cy!r}')
            print(f'     {n - diffs}/{n} rows identical')

    print('\n' + '=' * 78)
    print(f'TOTAL: {total_rows - total_diffs}/{total_rows} rows identical')
    if problems:
        print(f'GATE: FAIL — {len(problems)} problem(s)')
        for p in problems:
            print('  -', p)
        return 1
    print('GATE: PASS — every committed number reproduces from the raw dump')
    return 0

if __name__ == '__main__':
    sys.exit(main())
