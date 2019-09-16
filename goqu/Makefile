# Copyright 2011 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Author: conleyo@google.com (Conley Owens)

PKGSTEMS=quantum
EXAMPLESTEMS=deutsch deutsch-jozsa grover random shor simon

PKGDIRS=$(foreach stem, $(PKGSTEMS), src/$(stem))
EXAMPLEDIRS=$(foreach stem, $(EXAMPLESTEMS), examples/$(stem))
DIRS=$(PKGDIRS) $(EXAMPLEDIRS)

all:
	$(foreach dir, $(PKGDIRS), cd $(dir); make; cd -;)

install:
	$(foreach dir, $(PKGDIRS), cd $(dir); make install; cd -;)

test:
	$(foreach dir, $(PKGDIRS), cd $(dir); make test; cd -;)

all-examples:
	$(foreach dir, $(EXAMPLEDIRS), cd $(dir); make; cd -;)

clean:
	$(foreach dir, $(DIRS), cd $(dir); make clean; cd -;)
