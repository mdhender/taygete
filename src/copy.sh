#!/bin/bash

# taygete - a game engine for a game.
# Copyright (c) 2026 Michael D Henderson.
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

for file in	Makefile		\
		code.h			\
		dir.h			\
		display.h		\
		etc.h			\
		garr.h			\
		loc.h			\
		loop.h			\
		oly.h			\
		order.h			\
		sout.h			\
		stack.h			\
		swear.h			\
		u.h			\
		use.h			\
		z.h			\
		add.c			\
		adv.c			\
		alchem.c		\
		art.c			\
		basic.c			\
		beast.c			\
		build.c			\
		buy.c			\
		c1.c			\
		c2.c			\
		check.c			\
		cloud.c			\
		code.c			\
		combat.c		\
		day.c			\
		dir.c			\
		display.c		\
		eat.c			\
		faery.c			\
		garr.c			\
		gate.c			\
		glob.c			\
		gm.c			\
		hades.c			\
		immed.c			\
		input.c			\
		io.c			\
		loc.c			\
		lore.c			\
		main.c			\
		make.c			\
		move.c			\
		necro.c			\
		npc.c			\
		order.c			\
		perm.c			\
		produce.c		\
		pw.c			\
		quest.c			\
		relig.c			\
		report.c		\
		rnd.c			\
		savage.c		\
		scry.c			\
		seed.c			\
		sout.c			\
		stack.c			\
		stealth.c		\
		storm.c			\
		summary.c		\
		swear.c			\
		tunnel.c		\
		u.c			\
		use.c			\
		z.c
do
  cp -p ../olympia/${file} . || exit 2
done
