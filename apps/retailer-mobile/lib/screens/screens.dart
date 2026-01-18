// lib/screens/home_screen.dart
// Main dashboard for gig workers showing status, earnings, and tasks

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';

class HomeScreen extends ConsumerStatefulWidget {
  const HomeScreen({super.key});

  @override
  ConsumerState<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends ConsumerState<HomeScreen> {
  @override
  void initState() {
    super.initState();
    // Fetch initial data
    Future.microtask(() {
      ref.read(tasksProvider.notifier).fetchTasks();
    });
  }

  @override
  Widget build(BuildContext context) {
    final auth = ref.watch(authStateProvider);
    final availability = ref.watch(availabilityProvider);
    final tasks = ref.watch(tasksProvider);
    final location = ref.watch(locationProvider);

    final worker = auth.worker;
    final isOnline = availability.status == AvailabilityStatus.online;

    return Scaffold(
      body: SafeArea(
        child: RefreshIndicator(
          onRefresh: () async {
            await ref.read(tasksProvider.notifier).fetchTasks();
          },
          child: CustomScrollView(
            slivers: [
              // App Bar with greeting
              SliverAppBar(
                expandedHeight: 120,
                floating: true,
                pinned: true,
                flexibleSpace: FlexibleSpaceBar(
                  background: Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      gradient: LinearGradient(
                        colors: [
                          Theme.of(context).colorScheme.primary,
                          Theme.of(context).colorScheme.primaryContainer,
                        ],
                      ),
                    ),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.end,
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          _getGreeting(),
                          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                            color: Colors.white70,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          worker?.firstName ?? 'Worker',
                          style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                            color: Colors.white,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
                actions: [
                  // Notifications
                  IconButton(
                    icon: const Badge(
                      label: Text('3'),
                      child: Icon(Icons.notifications_outlined),
                    ),
                    onPressed: () {
                      // Show notifications
                    },
                  ),
                  // Profile
                  Padding(
                    padding: const EdgeInsets.only(right: 8),
                    child: CircleAvatar(
                      radius: 18,
                      backgroundColor: Colors.white24,
                      child: Text(
                        '${worker?.firstName[0] ?? ''}${worker?.lastName[0] ?? ''}',
                        style: const TextStyle(color: Colors.white),
                      ),
                    ),
                  ),
                ],
              ),

              // Online/Offline Toggle
              SliverToBoxAdapter(
                child: Container(
                  margin: const EdgeInsets.all(16),
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: isOnline 
                        ? Colors.green.withOpacity(0.1)
                        : Colors.grey.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(
                      color: isOnline ? Colors.green : Colors.grey,
                      width: 2,
                    ),
                  ),
                  child: Row(
                    children: [
                      // Status indicator
                      Container(
                        width: 48,
                        height: 48,
                        decoration: BoxDecoration(
                          color: isOnline ? Colors.green : Colors.grey,
                          shape: BoxShape.circle,
                        ),
                        child: Icon(
                          isOnline ? Icons.power_settings_new : Icons.power_off,
                          color: Colors.white,
                        ),
                      ),
                      const SizedBox(width: 16),
                      // Status text
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              isOnline ? 'You\'re Online' : 'You\'re Offline',
                              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                                fontWeight: FontWeight.bold,
                                color: isOnline ? Colors.green : Colors.grey,
                              ),
                            ),
                            if (isOnline && availability.onlineSince != null)
                              Text(
                                'Since ${DateFormat.jm().format(availability.onlineSince!)}',
                                style: Theme.of(context).textTheme.bodySmall,
                              )
                            else
                              const Text('Tap to start accepting tasks'),
                          ],
                        ),
                      ),
                      // Toggle switch
                      Switch(
                        value: isOnline,
                        onChanged: availability.isLoading
                            ? null
                            : (value) {
                                if (value) {
                                  ref.read(availabilityProvider.notifier).goOnline();
                                } else {
                                  _showGoOfflineDialog();
                                }
                              },
                        activeColor: Colors.green,
                      ),
                    ],
                  ),
                ),
              ),

              // Current Task Card
              if (tasks.currentTask != null)
                SliverToBoxAdapter(
                  child: _CurrentTaskCard(task: tasks.currentTask!),
                ),

              // Today's Stats
              SliverToBoxAdapter(
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 16),
                  child: _TodayStatsCard(
                    completedTasks: tasks.completedToday.length,
                    worker: worker,
                  ),
                ),
              ),

              // Section Header - Pending Tasks
              if (tasks.pendingTasks.isNotEmpty)
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.fromLTRB(16, 24, 16, 8),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          'Available Tasks',
                          style: Theme.of(context).textTheme.titleMedium?.copyWith(
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        TextButton(
                          onPressed: () => context.push('/tasks'),
                          child: const Text('View All'),
                        ),
                      ],
                    ),
                  ),
                ),

              // Task Offers List
              if (tasks.pendingTasks.isNotEmpty)
                SliverList(
                  delegate: SliverChildBuilderDelegate(
                    (context, index) {
                      final task = tasks.pendingTasks[index];
                      return _TaskOfferCard(
                        task: task,
                        onAccept: () => _acceptTask(task),
                        onDecline: () => _declineTask(task),
                      );
                    },
                    childCount: tasks.pendingTasks.take(3).length,
                  ),
                ),

              // Empty state when no tasks
              if (isOnline && tasks.pendingTasks.isEmpty && tasks.currentTask == null)
                SliverFillRemaining(
                  hasScrollBody: false,
                  child: Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(
                          Icons.hourglass_empty,
                          size: 64,
                          color: Theme.of(context).colorScheme.outline,
                        ),
                        const SizedBox(height: 16),
                        Text(
                          'Waiting for tasks...',
                          style: Theme.of(context).textTheme.titleMedium,
                        ),
                        const SizedBox(height: 8),
                        Text(
                          'Stay in a busy area to get more offers',
                          style: Theme.of(context).textTheme.bodySmall,
                        ),
                      ],
                    ),
                  ),
                ),

              // Bottom padding
              const SliverPadding(padding: EdgeInsets.only(bottom: 100)),
            ],
          ),
        ),
      ),
    );
  }

  String _getGreeting() {
    final hour = DateTime.now().hour;
    if (hour < 12) return 'Good Morning';
    if (hour < 17) return 'Good Afternoon';
    return 'Good Evening';
  }

  void _showGoOfflineDialog() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Go Offline?'),
        content: const Text(
          'You will stop receiving new task offers. '
          'Any accepted tasks will still need to be completed.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              ref.read(availabilityProvider.notifier).goOffline();
            },
            child: const Text('Go Offline'),
          ),
        ],
      ),
    );
  }

  void _acceptTask(Task task) {
    ref.read(tasksProvider.notifier).acceptTask(task.id);
  }

  void _declineTask(Task task) {
    showDialog(
      context: context,
      builder: (context) => _DeclineReasonDialog(
        onDecline: (reason) {
          Navigator.pop(context);
          ref.read(tasksProvider.notifier).declineTask(task.id, reason: reason);
        },
      ),
    );
  }
}

// Current Task Card Widget
class _CurrentTaskCard extends StatelessWidget {
  final Task task;

  const _CurrentTaskCard({required this.task});

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      color: Theme.of(context).colorScheme.primaryContainer,
      child: InkWell(
        onTap: () => context.push('/tasks/${task.id}'),
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: _getStatusColor(task.status),
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: Text(
                      _getStatusText(task.status),
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 12,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  const Spacer(),
                  Text(
                    '₦${task.earning?.total.toStringAsFixed(0) ?? '0'}',
                    style: Theme.of(context).textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 12),
              Row(
                children: [
                  const Icon(Icons.location_on, size: 16, color: Colors.red),
                  const SizedBox(width: 4),
                  Expanded(
                    child: Text(
                      task.dropoff.addressLine1,
                      style: Theme.of(context).textTheme.bodyMedium,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 8),
              Row(
                children: [
                  const Icon(Icons.person, size: 16),
                  const SizedBox(width: 4),
                  Text(task.customerName),
                  const Spacer(),
                  IconButton(
                    icon: const Icon(Icons.phone),
                    onPressed: () {
                      // Call customer
                    },
                  ),
                  IconButton(
                    icon: const Icon(Icons.navigation),
                    onPressed: () {
                      context.push('/tasks/${task.id}/navigate');
                    },
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Color _getStatusColor(String status) {
    switch (status) {
      case 'accepted':
        return Colors.blue;
      case 'en_route':
      case 'at_pickup':
        return Colors.orange;
      case 'picked_up':
      case 'in_transit':
        return Colors.purple;
      case 'at_dropoff':
        return Colors.green;
      default:
        return Colors.grey;
    }
  }

  String _getStatusText(String status) {
    switch (status) {
      case 'accepted':
        return 'ACCEPTED';
      case 'en_route':
        return 'EN ROUTE';
      case 'at_pickup':
        return 'AT PICKUP';
      case 'picked_up':
        return 'PICKED UP';
      case 'in_transit':
        return 'IN TRANSIT';
      case 'at_dropoff':
        return 'AT DROPOFF';
      default:
        return status.toUpperCase();
    }
  }
}

// Today's Stats Card
class _TodayStatsCard extends StatelessWidget {
  final int completedTasks;
  final Worker? worker;

  const _TodayStatsCard({
    required this.completedTasks,
    this.worker,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              "Today's Performance",
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: _StatItem(
                    icon: Icons.check_circle,
                    label: 'Completed',
                    value: '$completedTasks',
                    color: Colors.green,
                  ),
                ),
                Expanded(
                  child: _StatItem(
                    icon: Icons.star,
                    label: 'Rating',
                    value: '${worker?.rating.toStringAsFixed(1) ?? '0.0'}',
                    color: Colors.amber,
                  ),
                ),
                Expanded(
                  child: _StatItem(
                    icon: Icons.emoji_events,
                    label: 'Level',
                    value: '${worker?.level ?? 1}',
                    color: Colors.purple,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            // XP Progress
            if (worker != null) ...[
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    '${worker!.xp} XP',
                    style: Theme.of(context).textTheme.bodySmall,
                  ),
                  Text(
                    '${_xpForNextLevel(worker!.level)} XP for Level ${worker!.level + 1}',
                    style: Theme.of(context).textTheme.bodySmall,
                  ),
                ],
              ),
              const SizedBox(height: 4),
              LinearProgressIndicator(
                value: worker!.xp / _xpForNextLevel(worker!.level),
                backgroundColor: Colors.grey.shade200,
              ),
            ],
          ],
        ),
      ),
    );
  }

  int _xpForNextLevel(int currentLevel) {
    return currentLevel * 100 + 100; // Simple formula
  }
}

class _StatItem extends StatelessWidget {
  final IconData icon;
  final String label;
  final String value;
  final Color color;

  const _StatItem({
    required this.icon,
    required this.label,
    required this.value,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Icon(icon, color: color, size: 28),
        const SizedBox(height: 4),
        Text(
          value,
          style: Theme.of(context).textTheme.titleLarge?.copyWith(
            fontWeight: FontWeight.bold,
          ),
        ),
        Text(
          label,
          style: Theme.of(context).textTheme.bodySmall,
        ),
      ],
    );
  }
}

// Task Offer Card
class _TaskOfferCard extends StatelessWidget {
  final Task task;
  final VoidCallback onAccept;
  final VoidCallback onDecline;

  const _TaskOfferCard({
    required this.task,
    required this.onAccept,
    required this.onDecline,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Header with type and earning
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: Theme.of(context).colorScheme.secondaryContainer,
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Text(
                    task.type.toUpperCase(),
                    style: TextStyle(
                      color: Theme.of(context).colorScheme.onSecondaryContainer,
                      fontSize: 12,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
                if (task.priority == 'urgent' || task.priority == 'express')
                  Container(
                    margin: const EdgeInsets.only(left: 8),
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: Colors.red,
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: Text(
                      task.priority.toUpperCase(),
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 12,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                const Spacer(),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text(
                      '₦${task.earning?.total.toStringAsFixed(0) ?? '0'}',
                      style: Theme.of(context).textTheme.titleLarge?.copyWith(
                        fontWeight: FontWeight.bold,
                        color: Colors.green,
                      ),
                    ),
                    if (task.collectionAmount > 0)
                      Text(
                        'Collect: ₦${task.collectionAmount.toStringAsFixed(0)}',
                        style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: Colors.orange,
                        ),
                      ),
                  ],
                ),
              ],
            ),
            const Divider(height: 24),
            
            // Pickup location (if applicable)
            if (task.pickup != null) ...[
              Row(
                children: [
                  Container(
                    width: 24,
                    height: 24,
                    decoration: BoxDecoration(
                      color: Colors.blue.shade100,
                      shape: BoxShape.circle,
                    ),
                    child: const Icon(Icons.arrow_upward, size: 14, color: Colors.blue),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Text('PICKUP', style: TextStyle(fontSize: 10, color: Colors.grey)),
                        Text(
                          task.pickup!.addressLine1,
                          style: Theme.of(context).textTheme.bodyMedium,
                        ),
                      ],
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 8),
              // Dotted line
              Padding(
                padding: const EdgeInsets.only(left: 11),
                child: Container(
                  width: 2,
                  height: 20,
                  color: Colors.grey.shade300,
                ),
              ),
              const SizedBox(height: 8),
            ],
            
            // Dropoff location
            Row(
              children: [
                Container(
                  width: 24,
                  height: 24,
                  decoration: BoxDecoration(
                    color: Colors.red.shade100,
                    shape: BoxShape.circle,
                  ),
                  child: const Icon(Icons.arrow_downward, size: 14, color: Colors.red),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('DROPOFF', style: TextStyle(fontSize: 10, color: Colors.grey)),
                      Text(
                        task.dropoff.addressLine1,
                        style: Theme.of(context).textTheme.bodyMedium,
                      ),
                      Text(
                        task.dropoff.city,
                        style: Theme.of(context).textTheme.bodySmall,
                      ),
                    ],
                  ),
                ),
              ],
            ),
            
            const SizedBox(height: 16),
            
            // Time window
            Row(
              children: [
                const Icon(Icons.access_time, size: 16, color: Colors.grey),
                const SizedBox(width: 4),
                Text(
                  '${DateFormat.jm().format(task.timeWindowStart)} - ${DateFormat.jm().format(task.timeWindowEnd)}',
                  style: Theme.of(context).textTheme.bodySmall,
                ),
                const Spacer(),
                if (task.items.isNotEmpty)
                  Row(
                    children: [
                      const Icon(Icons.inventory_2, size: 16, color: Colors.grey),
                      const SizedBox(width: 4),
                      Text(
                        '${task.items.length} items',
                        style: Theme.of(context).textTheme.bodySmall,
                      ),
                    ],
                  ),
              ],
            ),
            
            const SizedBox(height: 16),
            
            // Action buttons
            Row(
              children: [
                Expanded(
                  child: OutlinedButton(
                    onPressed: onDecline,
                    child: const Text('Decline'),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  flex: 2,
                  child: ElevatedButton(
                    onPressed: onAccept,
                    child: const Text('Accept'),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

// Decline Reason Dialog
class _DeclineReasonDialog extends StatefulWidget {
  final Function(String?) onDecline;

  const _DeclineReasonDialog({required this.onDecline});

  @override
  State<_DeclineReasonDialog> createState() => _DeclineReasonDialogState();
}

class _DeclineReasonDialogState extends State<_DeclineReasonDialog> {
  String? _selectedReason;

  final _reasons = [
    'Too far away',
    'Vehicle issue',
    'Personal emergency',
    'Area not safe',
    'Too heavy',
    'Other',
  ];

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Why are you declining?'),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        children: _reasons.map((reason) {
          return RadioListTile<String>(
            title: Text(reason),
            value: reason,
            groupValue: _selectedReason,
            onChanged: (value) => setState(() => _selectedReason = value),
          );
        }).toList(),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context),
          child: const Text('Cancel'),
        ),
        ElevatedButton(
          onPressed: () => widget.onDecline(_selectedReason),
          child: const Text('Decline'),
        ),
      ],
    );
  }
}

// lib/screens/task_completion_screen.dart
// Screen for completing a delivery with proof

class TaskCompletionScreen extends ConsumerStatefulWidget {
  final String taskId;

  const TaskCompletionScreen({super.key, required this.taskId});

  @override
  ConsumerState<TaskCompletionScreen> createState() => _TaskCompletionScreenState();
}

class _TaskCompletionScreenState extends ConsumerState<TaskCompletionScreen> {
  final _formKey = GlobalKey<FormState>();
  final _receiverNameController = TextEditingController();
  final _notesController = TextEditingController();
  
  List<String> _photos = [];
  String? _signatureBase64;
  double _collectedAmount = 0;
  String _paymentMethod = 'cash';
  bool _isLoading = false;

  @override
  Widget build(BuildContext context) {
    final tasks = ref.watch(tasksProvider);
    final task = tasks.currentTask;

    if (task == null) {
      return const Scaffold(
        body: Center(child: Text('Task not found')),
      );
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text('Complete Delivery'),
      ),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            // Customer info
            Card(
              child: ListTile(
                leading: const CircleAvatar(child: Icon(Icons.person)),
                title: Text(task.customerName),
                subtitle: Text(task.dropoff.addressLine1),
              ),
            ),
            const SizedBox(height: 24),

            // Proof of Delivery Photos
            Text(
              'Proof of Delivery Photos',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            _PhotoGrid(
              photos: _photos,
              onAddPhoto: _addPhoto,
              onRemovePhoto: (index) {
                setState(() => _photos.removeAt(index));
              },
            ),
            const SizedBox(height: 24),

            // Signature
            Text(
              'Customer Signature',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            _SignaturePad(
              signature: _signatureBase64,
              onSigned: (signature) {
                setState(() => _signatureBase64 = signature);
              },
              onClear: () {
                setState(() => _signatureBase64 = null);
              },
            ),
            const SizedBox(height: 24),

            // Receiver name
            TextFormField(
              controller: _receiverNameController,
              decoration: const InputDecoration(
                labelText: 'Receiver Name',
                border: OutlineInputBorder(),
                prefixIcon: Icon(Icons.person_outline),
              ),
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter receiver name';
                }
                return null;
              },
            ),
            const SizedBox(height: 24),

            // Payment Collection (if applicable)
            if (task.collectionAmount > 0) ...[
              Text(
                'Payment Collection',
                style: Theme.of(context).textTheme.titleMedium,
              ),
              const SizedBox(height: 8),
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          const Text('Amount to collect:'),
                          Text(
                            '₦${task.collectionAmount.toStringAsFixed(2)}',
                            style: Theme.of(context).textTheme.titleLarge?.copyWith(
                              fontWeight: FontWeight.bold,
                              color: Colors.green,
                            ),
                          ),
                        ],
                      ),
                      const Divider(),
                      TextFormField(
                        initialValue: task.collectionAmount.toStringAsFixed(2),
                        keyboardType: TextInputType.number,
                        decoration: const InputDecoration(
                          labelText: 'Amount Collected',
                          prefixText: '₦ ',
                          border: OutlineInputBorder(),
                        ),
                        onChanged: (value) {
                          _collectedAmount = double.tryParse(value) ?? 0;
                        },
                      ),
                      const SizedBox(height: 12),
                      DropdownButtonFormField<String>(
                        value: _paymentMethod,
                        decoration: const InputDecoration(
                          labelText: 'Payment Method',
                          border: OutlineInputBorder(),
                        ),
                        items: const [
                          DropdownMenuItem(value: 'cash', child: Text('Cash')),
                          DropdownMenuItem(value: 'pos', child: Text('POS')),
                          DropdownMenuItem(value: 'transfer', child: Text('Bank Transfer')),
                        ],
                        onChanged: (value) {
                          setState(() => _paymentMethod = value!);
                        },
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 24),
            ],

            // Notes
            TextFormField(
              controller: _notesController,
              maxLines: 3,
              decoration: const InputDecoration(
                labelText: 'Notes (Optional)',
                border: OutlineInputBorder(),
                hintText: 'Any additional notes about the delivery',
              ),
            ),
            const SizedBox(height: 32),

            // Complete button
            ElevatedButton(
              onPressed: _isLoading ? null : _completeDelivery,
              style: ElevatedButton.styleFrom(
                padding: const EdgeInsets.symmetric(vertical: 16),
                backgroundColor: Colors.green,
              ),
              child: _isLoading
                  ? const SizedBox(
                      height: 20,
                      width: 20,
                      child: CircularProgressIndicator(color: Colors.white),
                    )
                  : const Text(
                      'COMPLETE DELIVERY',
                      style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                        color: Colors.white,
                      ),
                    ),
            ),
            const SizedBox(height: 16),

            // Unable to deliver
            OutlinedButton(
              onPressed: () => _showFailureReasonDialog(),
              style: OutlinedButton.styleFrom(
                padding: const EdgeInsets.symmetric(vertical: 16),
                foregroundColor: Colors.red,
              ),
              child: const Text('Unable to Deliver'),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _addPhoto() async {
    // Use image_picker to take photo
    // final picker = ImagePicker();
    // final image = await picker.pickImage(source: ImageSource.camera);
    // if (image != null) {
    //   // Upload to server and get URL
    //   setState(() => _photos.add(imageUrl));
    // }
  }

  Future<void> _completeDelivery() async {
    if (!_formKey.currentState!.validate()) return;
    if (_photos.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please take at least one photo')),
      );
      return;
    }

    setState(() => _isLoading = true);

    try {
      final location = ref.read(locationProvider).currentLocation;
      
      await ref.read(tasksProvider.notifier).completeTask(
        widget.taskId,
        TaskCompletion(
          proofPhotos: _photos,
          signature: _signatureBase64,
          receiverName: _receiverNameController.text,
          location: location,
          collectedAmount: _collectedAmount,
          paymentMethod: _paymentMethod,
          notes: _notesController.text,
        ),
      );

      if (mounted) {
        // Show success
        showDialog(
          context: context,
          barrierDismissible: false,
          builder: (context) => _CompletionSuccessDialog(
            earning: ref.read(tasksProvider).completedToday.last.earning,
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e')),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }

  void _showFailureReasonDialog() {
    showDialog(
      context: context,
      builder: (context) => _FailureReasonDialog(
        onFail: (reasonCode, notes) async {
          Navigator.pop(context);
          await ref.read(tasksProvider.notifier).failTask(
            widget.taskId,
            reasonCode,
            notes: notes,
          );
          if (mounted) {
            context.go('/');
          }
        },
      ),
    );
  }
}

// Photo Grid Widget
class _PhotoGrid extends StatelessWidget {
  final List<String> photos;
  final VoidCallback onAddPhoto;
  final Function(int) onRemovePhoto;

  const _PhotoGrid({
    required this.photos,
    required this.onAddPhoto,
    required this.onRemovePhoto,
  });

  @override
  Widget build(BuildContext context) {
    return Wrap(
      spacing: 8,
      runSpacing: 8,
      children: [
        ...photos.asMap().entries.map((entry) {
          return Stack(
            children: [
              Container(
                width: 80,
                height: 80,
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(8),
                  image: DecorationImage(
                    image: NetworkImage(entry.value),
                    fit: BoxFit.cover,
                  ),
                ),
              ),
              Positioned(
                top: -8,
                right: -8,
                child: IconButton(
                  icon: const Icon(Icons.cancel, color: Colors.red),
                  onPressed: () => onRemovePhoto(entry.key),
                ),
              ),
            ],
          );
        }),
        if (photos.length < 4)
          GestureDetector(
            onTap: onAddPhoto,
            child: Container(
              width: 80,
              height: 80,
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(8),
                border: Border.all(color: Colors.grey),
              ),
              child: const Icon(Icons.add_a_photo, color: Colors.grey),
            ),
          ),
      ],
    );
  }
}

// Signature Pad Widget
class _SignaturePad extends StatelessWidget {
  final String? signature;
  final Function(String) onSigned;
  final VoidCallback onClear;

  const _SignaturePad({
    this.signature,
    required this.onSigned,
    required this.onClear,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 150,
      decoration: BoxDecoration(
        border: Border.all(color: Colors.grey),
        borderRadius: BorderRadius.circular(8),
      ),
      child: signature != null
          ? Stack(
              children: [
                // Display signature
                Center(child: Image.memory(base64Decode(signature!))),
                Positioned(
                  top: 4,
                  right: 4,
                  child: IconButton(
                    icon: const Icon(Icons.clear),
                    onPressed: onClear,
                  ),
                ),
              ],
            )
          : Center(
              child: TextButton.icon(
                icon: const Icon(Icons.edit),
                label: const Text('Tap to sign'),
                onPressed: () {
                  // Open signature pad
                },
              ),
            ),
    );
  }
}

// Completion Success Dialog
class _CompletionSuccessDialog extends StatelessWidget {
  final Earning? earning;

  const _CompletionSuccessDialog({this.earning});

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      content: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          const Icon(Icons.check_circle, color: Colors.green, size: 64),
          const SizedBox(height: 16),
          const Text(
            'Delivery Complete!',
            style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
          ),
          if (earning != null) ...[
            const SizedBox(height: 8),
            Text(
              '+₦${earning!.total.toStringAsFixed(2)}',
              style: const TextStyle(
                fontSize: 24,
                fontWeight: FontWeight.bold,
                color: Colors.green,
              ),
            ),
          ],
        ],
      ),
      actions: [
        ElevatedButton(
          onPressed: () {
            Navigator.pop(context);
            context.go('/');
          },
          child: const Text('Continue'),
        ),
      ],
    );
  }
}

// Failure Reason Dialog
class _FailureReasonDialog extends StatefulWidget {
  final Function(String, String?) onFail;

  const _FailureReasonDialog({required this.onFail});

  @override
  State<_FailureReasonDialog> createState() => _FailureReasonDialogState();
}

class _FailureReasonDialogState extends State<_FailureReasonDialog> {
  String? _selectedReason;
  final _notesController = TextEditingController();

  final _reasons = {
    'customer_unavailable': 'Customer unavailable',
    'wrong_address': 'Wrong address',
    'customer_refused': 'Customer refused delivery',
    'damaged_goods': 'Goods damaged',
    'access_issue': 'Cannot access location',
    'other': 'Other',
  };

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Reason for failure'),
      content: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          ..._reasons.entries.map((entry) {
            return RadioListTile<String>(
              title: Text(entry.value),
              value: entry.key,
              groupValue: _selectedReason,
              onChanged: (value) => setState(() => _selectedReason = value),
            );
          }),
          if (_selectedReason == 'other')
            TextField(
              controller: _notesController,
              decoration: const InputDecoration(
                labelText: 'Please specify',
              ),
            ),
        ],
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context),
          child: const Text('Cancel'),
        ),
        ElevatedButton(
          onPressed: _selectedReason == null
              ? null
              : () => widget.onFail(_selectedReason!, _notesController.text),
          style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
          child: const Text('Confirm'),
        ),
      ],
    );
  }
}
