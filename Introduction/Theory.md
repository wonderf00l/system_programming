

- ***Unix***
	1. Name of a really existed operating system created at 1969 by Dennis Ritchie and Kenneth Thompson in Bell Labs company
	2. A set of operating systems, implementing the same programming interface and having similar concepts
	3. A concept, idea, to which many operating systems stick, and simplified version of which serves as a great teaching example
	4. *former trademark of AT&T company, and now of Open Group organisation making standards.
- Nothing is real "***Unix***" except an old really existing "Unix" system.
- Linux, MacOS, Solaris, BSD - all are slightly (more or less) **different variations of kernels, which are Unix-like**.



![[Pasted image 20241027171846.png]]

![[Pasted image 20241027172836.png]]
Some examples of the dynamic linking: firstly, again mmap - it not only loads a memory segment from a disk or main memory, but that segment extends process' address space, makes it literally bigger in terms of valid memory addresses. You can specify if the segment is code, or just data. Secondly, dynamic linking is about PIC - Position Independent Code. This is when compiled code doesn't have absolute addresses in it. It means it can be loaded into any address range and still be working. Usually we create such code using -fPIC GCC compiler option.



![[Pasted image 20241027172902.png]]
- The closer is a ring to the center, the more privileged it is in terms of access to hardware and being allowed to do actions affecting the whole system. In the center an OS kernel works. The OS defines rules by which a program from an external ring can access resources belonging to the programs in an internal ring.
- Nowadays we use only 2 - kernel space, user space.
- There is a new trend about virtualization of everything, and a so called ring '-1' appeared, a ring where works a hypervisor - a program below the kernel and managing multiple isolated kernels on one machine
	- **Гипервизор** — это **программа или аппаратная схема, которая позволяет одновременно, параллельно запускать несколько операционных систем на одном и том же хост-компьютере**.

Гипервизор также обеспечивает
- изоляцию операционных систем друг от друга
- защиту и безопасность
- разделение ресурсов между различными запущенными ОС и управление ресурсами

Гипервизор предоставляет работающим под его управлением на одном хост-компьютере ОС **средства связи и взаимодействия между собой** (например, через обмен файлами или сетевые соединения) так, как если бы эти ОС выполнялись на разных физических компьютерах.[1](https://ru.wikipedia.org/wiki/%D0%93%D0%B8%D0%BF%D0%B5%D1%80%D0%B2%D0%B8%D0%B7%D0%BE%D1%80 "ru.wikipedia.org – Гипервизор — Википедия")

**Выделяют два основных типа гипервизоров**:
- **Гипервизоры I типа**. Работают на «голом железе» и не требуют установки какой-либо операционной системы в качестве прослойки. Примеры: VMware ESXi, Hyper-V, KVM
- **Гипервизоры II типа**. Запускаются поверх операционной системы, установленной на хосте. Примеры: VMware Workstation, Oracle VirtualBox.


![[Pasted image 20241027173641.png]]
![[Pasted image 20241027173657.png]]



![[Pasted image 20241027174055.png]]



![[Pasted image 20241027174622.png]]
- OS consisting only of a kernel is nothing. ***When it is said Linux, we usually mean some complete operating system like Ubuntu. But in fact Linux is just the kernel. Not an operating system***. To use a kernel you need an environment - compilers, text editors, console, libraries, kernel modules, drivers, etc

**Linux — это семейство операционных систем, работающих на базе одноимённого ядра**. Нет одной операционной системы Linux, как, например, Windows или MacOS. Есть множество дистрибутивов, выполняющих конкретные задачи. [1](https://blog.skillfactory.ru/glossary/linux/)

Ядро Linux **бесплатное и распространяется по лицензии open source**. Это значит, что каждый разработчик может взять ядро и настроить по своему вкусу: добавить модули и программы, нарисовать любой интерфейс, внедрить продвинутые алгоритмы защиты и так далее. [2](https://skillbox.ru/media/code/chto-takoe-linux-gayd-po-samoy-svobodnoy-operatsionnoy-sisteme/)


![[Pasted image 20241027174833.png]]
- GNU includes gcc compiler, text editor emacs, lots of utilities like ls, grep, make, ld etc.



![[Pasted image 20241027183407.png]]
- Microkernels appeared as a way to simplify OS development, increase durability. 
- When OS consists of independent modules, they can be evolved and updated separately. 
- When a subsystem is down, it does not pull down all other subsystems and the entire machine, on the contrary with a monolithic kernel.(что-то похожее на микросервисную архитектуру)
- With growth of hardware performance it was noticed that microkernel has a serious problem in its design about the whole concept of kernel in multiple processes. Subsystems were forced to interact with each other via IPC - Inter Process Communication channels with a throughput bottleneck right inside the kernel.

![[Pasted image 20241027183848.png]]
- Общение между модулями ядра(разными процессами) посредством сообщений(через IPC) накладывает существенный оверхед на работу всей системы
![[Pasted image 20241027183857.png]]

![[Pasted image 20241027184408.png]]
- В данном случае все модули ядра находятся в 1-ом адресном пространстве, работающее ядро - 1 работающий процесс
- В данном случае вызов 1-ого модуля другим осуществляется напрямую через вызов функции ядра в каком-то другом пакете
- Нет расходов на IPC, но такие ядра сложнее в разработке и поддержке
- User interacts with monolithic kernel via system calls. It is a set of functions, which kernel exposes to be called by user's applications, and around which there are various wrappers. For instance, 'printf' is a wrapper around syscall 'write'. Function 'write' is also a wrapper around system call 'write', because system call is much lower than functions. Almost always highlevel functions called by user's code, and which are called system calls, actually are just wrappers, probably one-liners. 
- Another one of the main OS tasks - hardware management. It works via interrupts - kernel communicates with devices using electrical signals, interrupting a current thread and invoking the interrupt handler. Under devices usually is meant anything physical - memory, keyboard, mouse, screen, disks, CPU etc. Via system calls a user is able to work with the hardware through the kernel.


![[Pasted image 20241027185626.png]]
- System calls are not a Linux kernel's unique feature. During Unix systems boom every OS had its own set of system calls. Their number was growing, old calls were being deprecated or changed. The same was happening with even mere library functions, which could have one interface on multiple systems, but diff in behaviour. Corporations and institutes decided to create standards so as to simplify their life and life of their clients, users.

![[Pasted image 20241027185658.png]]
- In 1988 Institute of Electrical and Electronics Engineers (IEEE) released first POSIX - Portable Operating System Interface. The standard defines which services should provide OS to be "POSIX compliant". 
- POSIX defines only an interface, not implementation, it does not make a difference between library functions and system calls. It is an OS' job to choose how to implement something - in userspace via a library function, or in kernel space via a system call. The only thing which matters is the interface and behavior.
- Единый стандарт API, определенные требования к поведению ОС, набор функций - нужный для унификации, чтобы писать как можно более платформонезависимый код, а не разрабатывать одно и то же приложение по 10 раз для каждой ОС

---
презентация: https://slides.com/gerold103/sysprog_eng0
